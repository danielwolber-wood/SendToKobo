curl -X POST http://localhost:14000/v1/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<h1>Test Title</h1><p>This is a test paragraph with <strong>bold text</strong> and <em>italic text</em>.</p><ul><li>List item 1</li><li>List item 2</li></ul>",
    "title": "Test Document"
  }' \
  --output converted_output.epub
