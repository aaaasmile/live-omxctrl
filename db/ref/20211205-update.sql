-- Update file for changing the db
-- Delete unsused tables
BEGIN TRANSACTION;
DROP TABLE IF EXISTS `Item`;
DROP TABLE IF EXISTS `Current`;

DROP TABLE IF EXISTS `PlaylistItem`;
CREATE TABLE IF NOT EXISTS `PlaylistItem` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`playlist_id`	INTEGER,
	`URI`	TEXT,
	`Description`	TEXT,
	`MetaTitle`	TEXT,
	`MetaFileType`	TEXT,
	`MetaAlbum`	TEXT,
	`MetaArtist`	TEXT
);

--FileOrFolder : 0 folder, 1 Music File
CREATE TABLE IF NOT EXISTS `MusicFile` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`Timestamp`	INTEGER,
	`URI`	TEXT,
	`Title`	TEXT,
	`Description`	TEXT,
	`DurationInSec`	INTEGER,
	`FileOrFolder`	INTEGER,
  `ParentFolder`	TEXT,
  `MetaAlbum`	TEXT,
	`MetaArtist`	TEXT
);

COMMIT;