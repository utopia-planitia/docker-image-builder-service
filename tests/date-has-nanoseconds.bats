
load test_helper

@test "date has nanoseconds" {
  DATE_N=$(date +%N)
  echo "DATE_N=$DATE_N"
  [ "$DATE_N" != "" ]
}
