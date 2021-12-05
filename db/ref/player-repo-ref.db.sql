BEGIN TRANSACTION;
DROP TABLE IF EXISTS "playsearch";
CREATE VIRTUAL TABLE playsearch USING fts5(playsrowid, text);
DROP TABLE IF EXISTS "playsearch_data";
CREATE TABLE IF NOT EXISTS "playsearch_data" (
	"id"	INTEGER,
	"block"	BLOB,
	PRIMARY KEY("id")
);
DROP TABLE IF EXISTS "playsearch_idx";
CREATE TABLE IF NOT EXISTS "playsearch_idx" (
	"segid"	,
	"term"	,
	"pgno"	,
	PRIMARY KEY("segid","term")
) WITHOUT ROWID;
DROP TABLE IF EXISTS "playsearch_content";
CREATE TABLE IF NOT EXISTS "playsearch_content" (
	"id"	INTEGER,
	"c0"	,
	"c1"	,
	PRIMARY KEY("id")
);
DROP TABLE IF EXISTS "playsearch_docsize";
CREATE TABLE IF NOT EXISTS "playsearch_docsize" (
	"id"	INTEGER,
	"sz"	BLOB,
	PRIMARY KEY("id")
);
DROP TABLE IF EXISTS "playsearch_config";
CREATE TABLE IF NOT EXISTS "playsearch_config" (
	"k"	,
	"v"	,
	PRIMARY KEY("k")
) WITHOUT ROWID;
DROP TABLE IF EXISTS "Playlist";
CREATE TABLE IF NOT EXISTS "Playlist" (
	"id"	INTEGER,
	"Name"	TEXT UNIQUE,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "Item";
CREATE TABLE IF NOT EXISTS "Item" (
	"id"	INTEGER,
	"URI"	TEXT,
	"Info"	TEXT,
	"ItemType"	INTEGER,
	"Description"	TEXT,
	"MetaTitle"	TEXT,
	"MetaFileType"	TEXT,
	"MetaAlbum"	TEXT,
	"MetaArtist"	TEXT,
	"MetaAlbumArtist"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "History";
CREATE TABLE IF NOT EXISTS "History" (
	"id"	INTEGER,
	"Timestamp"	INTEGER,
	"URI"	TEXT,
	"Title"	TEXT,
	"Description"	TEXT,
	"Duration"	TEXT,
	"PlayPosition"	INTEGER,
	"DurationInSec"	INTEGER,
	"Type"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "Current";
CREATE TABLE IF NOT EXISTS "Current" (
	"id"	INTEGER,
	"ListName"	TEXT NOT NULL,
	"Volatile"	INTEGER,
	"URI"	TEXT,
	"Info"	TEXT,
	"ItemType"	INTEGER,
	PRIMARY KEY("id")
);
DROP TABLE IF EXISTS "Video";
CREATE TABLE IF NOT EXISTS "Video" (
	"id"	INTEGER,
	"Timestamp"	INTEGER,
	"URI"	TEXT,
	"Title"	TEXT,
	"Description"	TEXT,
	"Duration"	TEXT,
	"PlayPosition"	INTEGER,
	"DurationInSec"	INTEGER,
	"Type"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "PlaylistItem";
CREATE TABLE IF NOT EXISTS "PlaylistItem" (
	"id"	INTEGER,
	"playlist_id"	INTEGER,
	"item_id"	INTEGER,
	PRIMARY KEY("id")
);
DROP TABLE IF EXISTS "Radio";
CREATE TABLE IF NOT EXISTS "Radio" (
	"id"	INTEGER,
	"URI"	TEXT,
	"Name"	TEXT,
	"Description"	TEXT,
	"Genre"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
COMMIT;
