# encoding: utf-8
#
# Encoon : data structuration, presentation and navigation.
# 
# Copyright (C) 2012 David Lambert
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
# 
# See doc/COPYRIGHT.rdoc for more details.
for file in ["system", 
             "data_kinds", 
             "roles", 
             "system_administrator", 
             "countries", 
             "system_administror_credentials"]
  print "db:seeds: upload #{file}.xml"; STDOUT.flush
  upload = Upload.create
  upload.upload("db/#{file}.xml")
  puts "db:seeds: complete"
end