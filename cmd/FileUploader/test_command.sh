curl -X POST http://localhost/v1/api/upload \
  -F "token=@token.txt" \
  -F "filename=test.txt" \
  -F "filepath=/uploads" \
  -F "file=@test.txt"