puts "db:seeds: upload system.xml"
upload = Upload.create
upload.upload("db/system.xml")
puts "db:seeds: complete"