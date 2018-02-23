#!/bin/bash
#http://shout.setfive.com/2011/12/05/deleting-files-older-than-specified-time-with-s3cmd-and-bash/

olderThan=$((`date +%s` - $CACHE_MAX_AGE))

mc ls cache/$CACHE_BUCKET | while read -r line;
do
  createDate=`echo $line|awk {'print $1" "$2'}`
  createDate=`date -d"${createDate:1}" +%s`
  if [[ $createDate -lt $olderThan ]]; then 
      echo "deleteing $line"
      fileName=`echo $line|awk {'print $5'}`
      if [[ $fileName != "" ]]; then
          mc rm "cache/$CACHE_BUCKET/$fileName"
      fi
  fi
done;

mc ls -I cache/$CACHE_BUCKET | while read -r line;
do
  createDate=`echo $line|awk {'print $1" "$2'}`
  createDate=`date -d"${createDate:1}" +%s`
  if [[ $createDate -lt $olderThan ]]; then 
      echo "deleteing $line"
      fileName=`echo $line|awk {'print $5'}`
      if [[ $fileName != "" ]]; then
          mc rm "cache/$CACHE_BUCKET/$fileName"
      fi
  fi
done;
