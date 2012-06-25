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
class RowAttachment < ActiveRecord::Base
  belongs_to :row, 
             :foreign_key => "uuid", 
             :primary_key => "uuid"

  def attach_document=(input_file)
    log_debug "RowAttachment#document=(" + 
              "file_name=#{input_file.original_filename})"
    if input_file.present? 
      self.document = input_file.read
      self.content_type = input_file.content_type.chomp 
      self.file_name = input_file.original_filename 
    else
      self.document = nil
      self.content_type = nil 
      self.file_name = nil 
    end
  end
  
  def photo?
    self.document.present? and 
    self.content_type =~ /image+/
  end

  def document?
    self.document.present? and 
    not(self.content_type =~ /image+/)
  end
end