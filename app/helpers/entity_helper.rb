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
module EntityHelper
  def show_header_label(kind, label, description, list=false)
    if list
      case kind
        when COLUMN_TYPE_INTEGER then style = "list-header-number"
        when COLUMN_TYPE_DECIMAL then style = "list-header-number"
        when COLUMN_TYPE_DATE then style = "list-header-date"
        when COLUMN_TYPE_BOOLEAN then style = "list-header-boolean"
        else style = "list-header-string"
      end
    else
      style = "header"
    end
    content_tag(list ? "th" : "td",
      "<span>#{label}</span>".html_safe,
      :class => style
     )
  end

  def show_generic(value, list=false)
    content_tag("td", value, :class => (list ? "list-" : "") + "string")
  end

  def show_date(value, list=false)
    content_tag("td", value.nil? ? nil : I18n.l(value), :class => (list ? "list-" : "") + "date")
  end

  def show_number(value, list=false)
    content_tag("td", value, :class => (list ? "list-" : "") + "number")
  end

  def show_boolean(value, list=false)
    content_tag("td", "", 
        :class => (list ? "list-" : "") + "boolean " +
          ((value == "1" or value == "t" or value == "true") ? 
           (list ? "list-ok" : "ok") : 
           ""))
  end

  def show_reference(value, grid_uuid, referenced_link, referenced_name, referenced_description, list=false)
    content_tag("td", 
                (referenced_link.blank? ? "" : referenced_link.html_safe) + 
                (referenced_description.blank? ? "" : 
                  " <small>#{truncate_html(referenced_description)}</small>".html_safe), 
                :class => (list ? "list-" : "") + "string")
  end

  def show_hyperlink(value, list=false)
    content_tag("td", content_tag("a", value, :href => value), :class => (list ? "list-" : "") + "string")
  end

  def show_password(value, list=false)
    content_tag("td", icon('password'), :class => (list ? "list-" : "") + "string")
  end

  def show(value, kind, grid_uuid, referenced_link, referenced_name, referenced_description, list=false)
    case kind
      when COLUMN_TYPE_HYPERLINK then show_hyperlink(value, list)
      when COLUMN_TYPE_INTEGER then show_number(value, list)
      when COLUMN_TYPE_DECIMAL then show_number(value, list)
      when COLUMN_TYPE_DATE then show_date(value, list)
      when COLUMN_TYPE_BOOLEAN then show_boolean(value, list)
      when COLUMN_TYPE_REFERENCE then show_reference(value,
                                                 grid_uuid,
                                                 referenced_link,
                                                 referenced_name,
                                                 referenced_description,
                                                 list)
      when COLUMN_TYPE_PASSWORD then show_password(value, list)
      else show_generic(value, list)
    end
  end
  
  def show_entity(column, value, referenced_link, referenced_name, referenced_description)
    content_tag("tr",
      show_header_label(column.kind, 
                        column.name.html_safe, 
                        column.description.html_safe) +
      show(value, 
           column.kind, 
           column.grid_reference_uuid, 
           referenced_link, 
           referenced_name, 
           referenced_description)
    )
  end

  def show_entity_in_list(column, value, referenced_link, referenced_name, referenced_description = "")
    show(value,
         column.kind,
         column.grid_reference_uuid,
         referenced_link,
         referenced_name,
         referenced_description,
         true)
  end

  def edit_string(attribute, value)
    content_tag("td",
                text_field_tag(attribute,
                               value,
                               {:size => "70x1",
                                :id => "row_#{attribute}",
                                :name => "row_#{attribute}"}),
                :class => "string")
  end

  def edit_integer(attribute, value)
    content_tag("td",
                number_field_tag(attribute,
                                 value,
                                 {:size => "10x1",
                                  :id => "row_#{attribute}",
                                  :name => "row_#{attribute}"}),
                :class => "string")
  end

  def edit_text(attribute, value)
    content_tag("td",
                text_area_tag(attribute,
                               value,
                               {:size => "70x10",
                                :id => "row_#{attribute}",
                                :name => "row_#{attribute}"}),
                :class => "string")
  end

  def edit_boolean(attribute, value)
    content_tag("td",
                check_box_tag(attribute,
                              "t",
                              value,
                              {:size => "70x10",
                               :id => "row_#{attribute}",
                               :name => "row_#{attribute}"}),
                :class => "string")
  end

  def edit_date(attribute, value)
    content_tag("td",
                text_field_tag(attribute,
                               l(value),
                               {:size => "10x1",
                                :id => "row_#{attribute}",
                                :name => "row_#{attribute}",
                                :class => "datepicker"}),
                :class => "string")
  end

  def edit_password(attribute, value)
    content_tag("td",
                password_field_tag(attribute,
                                   value,
                                   {:size => "70x1",
                                    :id => "row_#{attribute}",
                                    :name => "row_#{attribute}"}),
                :class => "string")
  end

  def edit_reference(attribute, grid_uuid, include_blanks, value)
    collection = Grid.select_reference_rows(grid_uuid)
    unless collection.nil?
      content_tag("td", 
                  select_tag(attribute,
                             options_from_collection_for_select(collection, :uuid, :title, value),
                             {:include_blank => include_blanks,
                              :id => "row_#{attribute}",
                              :name => "row_#{attribute}"}),
                  :class => "string")
    end
  end

  def edit(attribute, value, kind, grid_uuid, include_blanks)
    case kind
      when COLUMN_TYPE_TEXT then edit_text(attribute, value)
      when COLUMN_TYPE_BOOLEAN then edit_boolean(attribute, value)
      when COLUMN_TYPE_INTEGER then edit_integer(attribute, value)
      when COLUMN_TYPE_DATE then edit_date(attribute, value)
      when COLUMN_TYPE_REFERENCE then edit_reference(attribute, grid_uuid, include_blanks, value)
      when COLUMN_TYPE_PASSWORD then edit_password(attribute, value)
      else edit_string(attribute, value)
    end
  end

  def show_edit_header_label(label, required, attribute)
    content_tag("td",
      content_tag("label",
                  label,
                  :for => "row_#{attribute}",
                  :class => (required ? "required" : "")),
      :class => "header"
    )
  end

  def edit_entity(row, column, description)
    attribute = row.initialized ? column.default_physical_column : column.physical_column
    value = row.read_value(column)
    content_tag("tr",
      show_edit_header_label(description.html_safe, column.required, attribute) +
      edit(attribute, value, column.kind, column.grid_reference_uuid, true),
      :id => "header-" + attribute
    )
  end
end