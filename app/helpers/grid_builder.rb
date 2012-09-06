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
class GridBuilder < ActionView::Helpers::FormBuilder
  def trunc(value, limit=150)
    return (value.length > limit) ? (value[0..limit-1] + "&hellip;") : value if value.present?
    ""
  end

  def edit_string(attribute)
    @template.content_tag("td", 
                          text_field(attribute, :size => "70x10"), 
                          :class => "string")
  end

  def edit_text(attribute)
    @template.content_tag("td", 
                          text_area(attribute, :size => "70x10"), 
                          :class => "string")
  end

  def edit_boolean(attribute)
    @template.content_tag("td", 
                          check_box(attribute, {}, "t", "f"), 
                          :class => "string")
  end

  def edit_date(attribute, value)
    @template.content_tag("td", 
                          text_field(attribute, :size => "10x1"),
                          :class => "string")
  end

  def edit_password(attribute)
    @template.content_tag("td", 
                          password_field(attribute, :size => "70x10"), 
                          :class => "string")
  end

  def edit_reference(attribute, grid_uuid, include_blanks)
    collection = Grid.select_reference_rows(grid_uuid)
    unless collection.nil?
      @template.content_tag("td", 
                            collection_select(attribute, 
                                              collection, 
                                              :uuid, 
                                              :title, 
                                              {:include_blank => include_blanks}), 
                            :class => "string")
    end
  end

  def edit(attribute, value, kind, grid_uuid, include_blanks)
    case kind
      when Column::STRING then edit_string(attribute)
      when Column::TEXT then edit_text(attribute)
      when Column::HYPERLINK then edit_string(attribute)
      when Column::BOOLEAN then edit_boolean(attribute)
      when Column::INTEGER then edit_string(attribute)
      when Column::DECIMAL then edit_string(attribute)
      when Column::DATE then edit_date(attribute, value)
      when Column::REFERENCE then edit_reference(attribute, grid_uuid, include_blanks)
      when Column::PASSWORD then edit_password(attribute)
      else ""
    end
  end

  def show_header_label(label, required, attribute)
    @template.content_tag("td",
      @template.content_tag("label", label, :for => attribute, :class => (required ? "required" : "")),
      :class => "header"
     )
  end

  def edit_entity(row, column, description)
    attribute = row.initialization? ? column.default_physical_column : column.physical_column
    value = row.read_value(column)
    @template.content_tag("tr",
      show_header_label(column.name.html_safe + description.html_safe, column.required, attribute) +
      edit(attribute, value, column.kind, column.grid_reference_uuid, true),
      :id => "header-" + attribute
    )
  end
end