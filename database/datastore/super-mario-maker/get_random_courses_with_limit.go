package datastore_smm_db

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"github.com/lib/pq"
)

// MaxRandomCourses is the most courses a single random course search will return.
// The largest request a real client makes is 100 (Course World on the Wii U)
const MaxRandomCourses = 100

const randomCourseColumns = `
	object.data_id,
	object.owner,
	object.size,
	object.name,
	object.data_type,
	object.meta_binary,
	object.permission,
	object.permission_recipients,
	object.delete_permission,
	object.delete_permission_recipients,
	object.period,
	object.refer_data_id,
	object.flag,
	object.tags,
	object.creation_date,
	object.update_date`

// * Picks courses by probing the custom rankings primary key
// * at random DataIDs, rather than sorting the whole table.
// * Each probe is a single index lookup which takes the first
// * available course at or after a random DataID, so the cost
// * depends on the number of probes and not the number of
// * courses. Multiple probes can land on the same course, so
// * more probes than needed are made ($1) and the duplicates
// * are dropped before taking the final amount ($2)
// *
// * Courses which come after a gap in DataIDs (deleted courses,
// * non-course objects, etc.) are slightly more likely to be
// * picked. This is fine for showing players random courses
const sampleRandomCoursesQuery = `
	WITH bounds AS (
		SELECT MIN(data_id) AS lowest, MAX(data_id) AS highest
		FROM datastore.object_custom_rankings
		WHERE application_id = 0
	),
	probes AS (
		SELECT (bounds.lowest + FLOOR(RANDOM() * (bounds.highest - bounds.lowest + 1)))::bigint AS data_id
		FROM bounds, GENERATE_SERIES(1, $1)
		WHERE bounds.lowest IS NOT NULL
	),
	picked AS (
		SELECT DISTINCT course.data_id, course.value
		FROM probes
		CROSS JOIN LATERAL (
			SELECT ranking.data_id, ranking.value
			FROM datastore.object_custom_rankings ranking
			JOIN datastore.objects object ON object.data_id = ranking.data_id
			WHERE
				ranking.data_id >= probes.data_id AND
				ranking.application_id = 0 AND
				object.upload_completed = TRUE AND
				object.deleted = FALSE AND
				object.under_review = FALSE
			ORDER BY ranking.data_id
			LIMIT 1
		) course
	)
	SELECT` + randomCourseColumns + `,
		picked.value
	FROM picked
	JOIN datastore.objects object ON object.data_id = picked.data_id
	ORDER BY RANDOM()
	LIMIT $2`

// * Sorts every available course. Only used when there are
// * too few courses for the sampling query to find enough
const allRandomCoursesQuery = `
	SELECT` + randomCourseColumns + `,
		ranking.value
	FROM datastore.objects object
	JOIN datastore.object_custom_rankings ranking
	ON
		object.data_id = ranking.data_id AND
		object.upload_completed = TRUE AND
		object.deleted = FALSE AND
		object.under_review = FALSE AND
		ranking.application_id = 0
	ORDER BY RANDOM()
	LIMIT $1`

func GetRandomCoursesWithLimit(limit int) (types.List[datastore_super_mario_maker_types.DataStoreCustomRankingResult], *nex.Error) {
	if limit <= 0 {
		return types.NewList[datastore_super_mario_maker_types.DataStoreCustomRankingResult](), nil
	}

	limit = min(limit, MaxRandomCourses)

	courses, nexError := queryRandomCourses(sampleRandomCoursesQuery, limit*2, limit)
	if nexError != nil {
		return nil, nexError
	}

	if len(courses) < limit {
		// * Sampling only comes up short when there are very few
		// * courses, where sorting all of them is cheap
		courses, nexError = queryRandomCourses(allRandomCoursesQuery, limit)
		if nexError != nil {
			return nil, nexError
		}
	}

	return attachCourseRatings(courses)
}

func queryRandomCourses(query string, args ...any) (types.List[datastore_super_mario_maker_types.DataStoreCustomRankingResult], *nex.Error) {
	courses := types.NewList[datastore_super_mario_maker_types.DataStoreCustomRankingResult]()

	rows, err := database.Postgres.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		course := datastore_super_mario_maker_types.NewDataStoreCustomRankingResult()
		course.Order = 0 // * Order is ALWAYS 0
		course.MetaInfo = datastore_types.NewDataStoreMetaInfo()
		course.MetaInfo.Permission = datastore_types.NewDataStorePermission()
		course.MetaInfo.DelPermission = datastore_types.NewDataStorePermission()
		course.MetaInfo.ExpireTime = types.NewDateTime(0x9C3F3E0000) // * 9999-12-31T00:00:00.000Z. This is what the real server sends
		course.MetaInfo.Ratings = types.NewList[datastore_types.DataStoreRatingInfoWithSlot]()

		var createdDate time.Time
		var updatedDate time.Time
		var tagArray []string

		err := rows.Scan(
			&course.MetaInfo.DataID,
			&course.MetaInfo.OwnerID,
			&course.MetaInfo.Size,
			&course.MetaInfo.Name,
			&course.MetaInfo.DataType,
			&course.MetaInfo.MetaBinary,
			&course.MetaInfo.Permission.Permission,
			pq.Array(&course.MetaInfo.Permission.RecipientIDs),
			&course.MetaInfo.DelPermission.Permission,
			pq.Array(&course.MetaInfo.DelPermission.RecipientIDs),
			&course.MetaInfo.Period,
			&course.MetaInfo.ReferDataID,
			&course.MetaInfo.Flag,
			pq.Array(&tagArray),
			&createdDate,
			&updatedDate,
			&course.Score,
		)
		if err != nil {
			globals.Logger.Error(err.Error())
			continue
		}

		course.MetaInfo.Tags = make(types.List[types.String], 0, len(tagArray))
		for i := range tagArray {
			course.MetaInfo.Tags = append(course.MetaInfo.Tags, types.String(tagArray[i]))
		}

		course.MetaInfo.CreatedTime.FromTimestamp(createdDate)
		course.MetaInfo.UpdatedTime.FromTimestamp(updatedDate)
		course.MetaInfo.ReferredTime.FromTimestamp(createdDate) // * This is what the real server does

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return courses, nil
}

func attachCourseRatings(courses types.List[datastore_super_mario_maker_types.DataStoreCustomRankingResult]) (types.List[datastore_super_mario_maker_types.DataStoreCustomRankingResult], *nex.Error) {
	dataIDs := make([]types.UInt64, 0, len(courses))
	for i := range courses {
		dataIDs = append(dataIDs, courses[i].MetaInfo.DataID)
	}

	ratings, nexError := datastore_db.GetObjectRatingsWithSlotByDataIDs(dataIDs)
	if nexError != nil {
		return nil, nexError
	}

	for i := range courses {
		if courseRatings, ok := ratings[courses[i].MetaInfo.DataID]; ok {
			courses[i].MetaInfo.Ratings = courseRatings
		}
	}

	return courses, nil
}
