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
    style = "header"
    if list
      case kind
        when Column::STRING then style = "list-header-string"
        when Column::TEXT then style = "list-header-string"
        when Column::HYPERLINK then style = "list-header-string"
        when Column::INTEGER then style = "list-header-number"
        when Column::DECIMAL then style = "list-header-number"
        when Column::DATE then style = "list-header-date"
        when Column::BOOLEAN then style = "list-header-boolean"
        when Column::REFERENCE then style = "list-header-string"
        when Column::PASSWORD then style = "list-header-string"
      end
    end
    content_tag(list ? "th" : "td",
      (("<span title=\"" + trunc(description) + "\">").html_safe +
      label +
      "</span>".html_safe),
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
    content_tag("td", (value == "1" or value == "t" or value == "true") ? 
        tag("img", { :src => asset_path("true.gif"), 
                     :height => "12", 
                     :width => "12" }, false, false).html_safe : "&nbsp;".html_safe, 
        :class => (list ? "list-" : "") + "boolean")
  end

  def show_reference(value, grid_uuid, referenced_link, referenced_name, referenced_description, list=false)
    content_tag("td", 
                (referenced_link.blank? ? "" : referenced_link.html_safe) + 
                (referenced_description.blank? ? "" : 
                  "<span class='small-description'>#{trunc(referenced_description)}</span>".html_safe), 
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
      when Column::STRING then show_generic(value, list)
      when Column::TEXT then show_generic(value, list)
      when Column::HYPERLINK then show_hyperlink(value, list)
      when Column::INTEGER then show_number(value, list)
      when Column::DECIMAL then show_number(value, list)
      when Column::DATE then show_date(value, list)
      when Column::BOOLEAN then show_boolean(value, list)
      when Column::REFERENCE then show_reference(value, grid_uuid, referenced_link, referenced_name, referenced_description, list)
      when Column::PASSWORD then show_password(value, list)
      else ""
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
    show(value, column.kind, column.grid_reference_uuid, referenced_link, referenced_name, referenced_description, true)
  end
end