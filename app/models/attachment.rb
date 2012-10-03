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
class Attachment < ActiveRecord::Base
  belongs_to :row, :foreign_key => "uuid", :primary_key => "uuid"

  # Paperclip interpolation rule: used to include row uuid in attachment paths.
  Paperclip.interpolates :uuid do |attachment, style|
    attachment.instance.uuid + "/" + attachment.original_filename
  end
  
  # Attached documents using PaperClip.
  has_attached_file :document, 
                    :styles => lambda { |attachment| !attachment.instance.photo? ? {} : { 
                       :inline => attachment.instance.reduce(580),
                       :mini => attachment.instance.crop(150,150),
                       :icon => attachment.instance.crop(32,32),
                       :micro => attachment.instance.crop(12,12)
                      } 
                    }, 
                    :path => ":rails_root/public/system/:style/:uuid",
                    :url => "/system/:style/:uuid"

  # Reduces a thumbnail based on a maximum width, 
  # if the original width is greater than this maximum.
  def reduce max_width
    if photo?
      geo = Paperclip::Geometry.from_file(document.to_file(:original))
      geo.width > max_width ? "#{max_width}" : "#{geo.width}"
    else
      nil
    end
  end

  # Resizes a thumbnail based on its orientation. 
  # Landscape pictures are cropped to a maximum width and height.
  # Portrait pictures are reduced to a maximum width and height, pictures are centered.
  def crop max_width, max_height
    if photo?
      geo = Paperclip::Geometry.from_file(document.to_file(:original))
      geo.width > geo.height ? "#{max_width}x#{max_height}#" : "#{max_width}x#{max_height}>"
    else
      nil
    end
  end
  
  def photo?
    self.document_content_type =~ /image+/
  end

  def document?
    not(self.document_content_type =~ /image+/)
  end
end