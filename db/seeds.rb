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
errors = total_count = total_inserted = total_updated = total_skipped = total_elapsed = 0
for file in ["administrator",
             "system_credentials",
             "system",
             "kinds",
             "display-modes",
             "sort-orders",
             "roles",
             "security_credentials",
             "security",
             "home_credentials",
             "home",
             "pages",
             "countries"]
  Entity.session_locale = 'en'
  Entity.session_as_of_date = Date.current
  Entity.session_user_uuid = SYSTEM_ADMINISTRATOR_UUID
  begin
    upload = Upload.create
    upload.upload("db/#{file}.xml")
    upload.save!
    Entity.log_debug "rake db:seed #{file} " +
                     "(#{upload.records} records): " + 
                     "#{upload.inserted} inserted, " + 
                     "#{upload.updated} updated, " +
                     "#{upload.skipped} skipped " +
                     "(#{upload.elapsed} ms).", true
    total_count = total_count + upload.records
    total_inserted = total_inserted + upload.inserted
    total_updated = total_updated + upload.updated
    total_skipped = total_skipped + upload.skipped
    total_elapsed = total_elapsed + upload.elapsed
  rescue Exception => invalid
    Entity.log_error "rake db:seed", invalid
    puts "rake db:seed " + invalid.inspect
    errors = errors + 1
  end
end
Entity.log_debug "rake db:seed total " +
                 "#{total_count} records, " + 
                 "#{errors} errors, " + 
                 "#{total_inserted} inserted, " + 
                 "#{total_updated} updated, " +
                 "#{total_skipped} skipped " +
                 "(#{total_elapsed} ms).", true