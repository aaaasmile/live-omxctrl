#!/bin/bash
for file in ./*.zip ; do
	fname=$(basename "$file")
	#echo "Name is $fname"
done

read -p "Copy the the live-omxctrl package $fname (y/n)? " -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Nn]$ ]]
then
	echo "Copy canceled"
	exit 0
fi

echo "Start to upload $fname"
rsync -avz $fname igors@pi4:/home/igors/app/live-omxctrl/zips

echo "That's all folks!"
