import dotenv from "dotenv";
import { input, confirm } from "@inquirer/prompts";
import pg from "pg";

dotenv.config({
  quiet: true
});

const client = new pg.Client({
  connectionString: process.env.POSTGRES_URL,
});

async function main() {
  await client.connect();

  const dataId = await input({ message: "Enter data_id:" });

  // Fetch object
  const queryResult = await client.query(
    "SELECT name, deleted FROM datastore.objects WHERE data_id = $1",
    [dataId]
  );

  if (queryResult.rows.length === 0) {
    console.log("❌ No object found with that id.");
    await client.end();
    return;
  }

  const name = queryResult.rows[0].name;
  const deleted = queryResult.rows[0].deleted;
  console.log(`Found object: "${name}"`);

  if (deleted) {
    console.log("❎ Already deleted.");
    await client.end();
    return;
  }

  const shouldDelete = await confirm({
    message: "Mark this object as deleted?",
    default: false,
  });

  if (shouldDelete) {
    await client.query(
      "UPDATE datastore.objects SET deleted = TRUE WHERE data_id = $1",
      [dataId]
    );
    console.log("✅ Object marked as deleted.");
  } else {
    console.log("❎ Cancelled.");
  }

  await client.end();
}

main().catch(async (err) => {
  console.error(err);
  await client.end();
});
