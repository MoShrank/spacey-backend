// inserts test user with following credentials:
// email: admin@spacey-learn.com
// password: 123456

// --> betaUser will be false because migration overrides field

db = db.getSiblingDB("spacey");
db.createCollection("user");

db.user.insertOne({
  _id: ObjectId("62fe2c8a657d9640f438717e"),
  name: "admin",
  email: "admin@spacey-learn.com",
  password: "$2a$14$czLG9a8oYcSSOqAXo0e.WeBz/qwrFLuK1qd3HTZHpBh6EwkDV.w.6",
  betaUser: true,
});
