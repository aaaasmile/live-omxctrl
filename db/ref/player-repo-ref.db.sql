BEGIN TRANSACTION;
DROP TABLE IF EXISTS "MusicFile";
CREATE TABLE IF NOT EXISTS "MusicFile" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Timestamp"	INTEGER,
	"URI"	TEXT,
	"Title"	TEXT,
	"Description"	TEXT,
	"DurationInSec"	INTEGER,
	"FileOrFolder"	INTEGER,
	"ParentFolder"	TEXT,
	"MetaAlbum"	TEXT,
	"MetaArtist"	TEXT
);
DROP TABLE IF EXISTS "PlaylistItem";
CREATE TABLE IF NOT EXISTS "PlaylistItem" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"playlist_id"	INTEGER,
	"URI"	TEXT,
	"Description"	TEXT,
	"MetaTitle"	TEXT,
	"MetaFileType"	TEXT,
	"MetaAlbum"	TEXT,
	"MetaArtist"	TEXT
);
DROP TABLE IF EXISTS "History";
CREATE TABLE IF NOT EXISTS "History" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Timestamp"	INTEGER,
	"URI"	TEXT,
	"Title"	TEXT,
	"Description"	TEXT,
	"Duration"	TEXT,
	"PlayPosition"	INTEGER,
	"DurationInSec"	INTEGER,
	"Type"	TEXT
);
DROP TABLE IF EXISTS "Playlist";
CREATE TABLE IF NOT EXISTS "Playlist" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Name"	TEXT UNIQUE
);
DROP TABLE IF EXISTS "Radio";
CREATE TABLE IF NOT EXISTS "Radio" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"URI"	TEXT,
	"Name"	TEXT,
	"Description"	TEXT,
	"Genre"	INTEGER
);
DROP TABLE IF EXISTS "Video";
CREATE TABLE IF NOT EXISTS "Video" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Timestamp"	INTEGER,
	"URI"	TEXT,
	"Title"	TEXT,
	"Description"	TEXT,
	"Duration"	TEXT,
	"PlayPosition"	INTEGER,
	"DurationInSec"	INTEGER,
	"Type"	TEXT
);
DROP TABLE IF EXISTS "playsearch_config";
DROP TABLE IF EXISTS "playsearch_docsize";
DROP TABLE IF EXISTS "playsearch_content";
DROP TABLE IF EXISTS "playsearch_idx";
DROP TABLE IF EXISTS "playsearch_data";
DROP TABLE IF EXISTS "playsearch";
-- Remember that CREATE playsearch generate automatically all other tables
CREATE VIRTUAL TABLE playsearch USING fts5(playsrowid, text);
COMMIT;
