curl -X POST http://localhost/v1/api/extract \
  -H "Content-Type: text/html" \
  -d '<html><head><title>Test Article</title></head><body><h1>Main Heading</h1><p>This is a test paragraph with some content to extract.</p><p>Another paragraph with more text content.</p></body></html>'