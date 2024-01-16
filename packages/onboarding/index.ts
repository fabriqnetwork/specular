import main from "./src";

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
