const app = new Application();
const router = new Router();

// Connect to MongoDB
const client = new MongoClient();
client.connectWithUri("mongodb://localhost:27017");
const db = client.database("meadRecipes");
const meadCollection = db.collection("recipes");

// Get all mead recipes
router.get("/mead", async (ctx) => {
  const recipes = await meadCollection.find();
  ctx.response.body = recipes;
});

// Get a specific mead recipe by ID
router.get("/mead/:id", async (ctx) => {
  const recipe = await meadCollection.findOne({ _id: { $oid: ctx.params.id } });
  if (recipe) {
    ctx.response.body = recipe;
  } else {
    ctx.response.status = 404;
    ctx.response.body = { message: "Recipe not found" };
  }
});

// Create a new mead recipe
router.post("/mead", async (ctx) => {
  const { name, ingredients, steps } = await ctx.request.body().value;
  const insertId = await meadCollection.insertOne({ name, ingredients, steps });
  ctx.response.status = 201;
  ctx.response.body = { id: insertId };
});

// Update a mead recipe by ID
router.put("/mead/:id", async (ctx) => {
  const { name, ingredients, steps } = await ctx.request.body().value;
  const result = await meadCollection.updateOne(
    { _id: { $oid: ctx.params.id } },
    { $set: { name, ingredients, steps } }
  );
  if (result.modifiedCount) {
    ctx.response.body = { message: "Recipe updated successfully" };
  } else {
    ctx.response.status = 404;
    ctx.response.body = { message: "Recipe not found" };
  }
});

// Delete a mead recipe by ID
router.delete("/mead/:id", async (ctx) => {
  const result = await meadCollection.deleteOne({ _id: { $oid: ctx.params.id } });
  if (result.deletedCount) {
    ctx.response.body = { message: "Recipe deleted successfully" };
  } else {
    ctx.response.status = 404;
    ctx.response.body = { message: "Recipe not found" };
  }
});

app.use(router.routes());
app.use(router.allowedMethods());

console.log("Server is running on http://localhost:8000");
await app.listen({ port: 8000 });

