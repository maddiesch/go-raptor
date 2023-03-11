--- This is the test file that will be applied to new test connections

CREATE TABLE "People" (
  "ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "FirstName" TEXT NOT NULL,
  "LastName" TEXT NOT NULL
);

CREATE UNIQUE INDEX "Index_People_Unique_Name" ON "People" ("FirstName", "LastName");

CREATE TABLE "Pets" (
  "ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "ParentID" INTEGER NOT NULL,
  "Type" TEXT NOT NULL,
  "Name" TEXT NOT NULL,
  "Age" INTEGER
);

INSERT INTO "People" ("FirstName", "LastName") VALUES ('Maddie', 'Schipper');
INSERT INTO "People" ("FirstName", "LastName") VALUES ('Elle', 'Woods');
INSERT INTO "People" ("FirstName", "LastName") VALUES ('Jackson', 'Briggs');

INSERT INTO "Pets" ("ParentID", "Type", "Name", "Age") VALUES ((SELECT "ID" FROM "People" WHERE "FirstName" = 'Maddie' AND "LastName" = 'Schipper'), 'Dog', 'Sterling', 5);
INSERT INTO "Pets" ("ParentID", "Type", "Name") VALUES ((SELECT "ID" FROM "People" WHERE "FirstName" = 'Elle' AND "LastName" = 'Woods'), 'Dog', 'Bruiser');
INSERT INTO "Pets" ("ParentID", "Type", "Name") VALUES ((SELECT "ID" FROM "People" WHERE "FirstName" = 'Jackson' AND "LastName" = 'Briggs'), 'Dog', 'Lulu');
