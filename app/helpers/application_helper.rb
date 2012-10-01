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
module ApplicationHelper
  # Defines special builder for data grid forms
  def grid_form_for(name, *args, &block)
    options = args.last.is_a?(Hash) ? args.pop : {}
    options = options.merge(:builder => GridBuilder)
    args = (args << options)
    form_for(name, *args, &block)
  end
  
  def current_page_title
    @page_title.present? ? (@page_title + " | " + current_application_title) : current_application_title 
  end
  
  def current_application_title
    APPLICATION_TITLE + ('production' != Rails.env ? (" (" + Rails.env + ")") : "")
  end
  
  # Displays an icon using its name
  def icon(name, title=nil)
    # Note: icons library: http://www.small-icons.com/packs/24x24-free-button-icons.jpg
    if name.present?
      tag("img", {:src => asset_path(name + ".gif"), 
                  :height => "12", 
                  :width => "12", 
                  :border => "0",
                  :title => title})
    else
      ""
    end
  end
  
  # Displays an icon using its name
  def bigicon(name)
    # Note: icons library: http://www.small-icons.com/packs/24x24-free-button-icons.jpg
    if name.present?
      tag("img", {:src => asset_path(name + ".gif"), 
                  :height => "18", 
                  :width => "18", 
                  :border => "0"})
    else
      ""
    end
  end
  
  # Displays an icon using its name
  def i(name)
    if name.present?
      tag("img", {:src => asset_path(name + ".png"), 
                  :height => "24", 
                  :width => "24", 
                  :border => "0"})
    else
      ""
    end
  end
  
  def display_filters(grid, filters, search=nil)
    output = ""
    if filters.present?
      filters.collect do |filter|
        column_uuid = filter[:column_uuid]
        row_name = filter[:row_name]
        column = grid.column_select_entity_by_uuid(column_uuid)
        if column.present? and column.kind == Column::REFERENCE
          output << ", " if output != ""
          output << "\"#{column.name}\" equals to \"#{row_name}\"".html_safe
        end
      end
    end
    if search.present?
      output << ", " if output != ""
      output << icon('yellowlight') + 
                "\"Name\" or \"description\" contains \"#{search}\"".html_safe
    end
    output
  end

  # This is used to generate a unique id in the page
  # in order to avoid duplicates.
  # This is based on the concatenation of identifiers used to filter data 
  def get_filters_uuid(filters)
    output = ""
    if filters.present?
      filters.collect do |filter|
         output << filter[:column_uuid]
      end
    end
    output
  end

  def display_history(history)
    output = ""
    if history.present?
      history.reverse_each do |link|
        hyperlink = content_tag("a", link[:page_title].html_safe, :href => link[:url])
        hyperlink = link[:page_title] if link[:url] == request.url
        output << content_tag("li", 
                                icon(link[:page_icon].to_s) + 
                                "&nbsp;".html_safe + 
                                hyperlink +
                                "&nbsp;".html_safe + 
                                content_tag("span", 
                                    t('general.ago', 
                                           :time => time_ago_in_words(link[:when], 
                                                     :include_seconds => true)), 
                                    :class => 'detail'), 
                                :class => "description")
      end
    end
    output
  end
  
  def display_date(date)
    (date.blank? or date == Entity.begin_of_time or date == Entity.end_of_time) ? nil : date
  end

  def display_begin_date(date)
    date == Entity.begin_of_time ? t('general.undefined') : date
  end

  def display_end_date(date)
    date == Entity.end_of_time ? t('general.undefined') : date
  end

  def display_distance_date(date)
    if date != Entity.begin_of_time and date != Entity.end_of_time
      now = Time.now
      today = Date::civil(now.year, now.month, now.day)
      if date > today
        t('general.ahead', :time => time_ago_in_words(date))
      else
        t('general.ago', :time => time_ago_in_words(date))
      end
    end
  end
  
  def display_new
    content_tag("div", t('general.new'), :class => "new")
  end
    
  def display_updated_date(entity)
    t(entity.revision == 1 ? 'general.created' : 'general.updated', 
            :time => time_ago_in_words(entity.updated_at, :include_seconds => true)) +
    ((Time.now-entity.updated_at) < 1.day ? display_new : "")
  end
  
  def display_created_time_by(entity, who, who_uuid)
    t('general.time_by', 
            :time => time_ago_in_words(
                                entity.created_at, 
                                :include_seconds => true),
            :by => (who.present? and who_uuid.present?) ? 
                       link_to_unless_current(who, 
                           show_path(:workspace => Workspace::SYSTEM_WORKSPACE_URI,
                                     :grid => User::ROOT_UUID,
                                     :row => who_uuid)) :
                       t('general.unknown'))
  end

  def display_updated_time_by(entity, who, who_uuid)
    t('general.time_by', 
            :time => time_ago_in_words(
                                entity.updated_at, 
                                :include_seconds => true),
            :by => (who.present? and who_uuid.present?) ? 
                       link_to_unless_current(who, 
                           show_path(:workspace => Workspace::SYSTEM_WORKSPACE_URI,
                                     :grid => User::ROOT_UUID,
                                     :row => who_uuid)) :
                       t('general.unknown')) +
    ((Time.now-entity.updated_at) < 1.day ? display_new : "")
  end

  def display_warning_current_date(as_of_date)
    now = Time.now
    today = Date::civil(now.year, now.month, now.day)
    as_of_date != today ? icon('warning') : ""
  end
  
  def information(entity, show_required=false)
    output = ""
    if entity.begin != Entity.begin_of_time
      output << t('general.begins', :time => entity.begin.to_s)
    end
    if entity.end != Entity.end_of_time
      output << " | " if output.length>0
      output <<  t('general.ends', :time => entity.end.to_s)
    end
    if show_required
      output << " | " if output.length>0 and show_required
      output << t('general.version', :version => entity.version.to_s)
    end
    if not entity.enabled
      output << " " if output.length>0
      output << t('general.inactive')
      output << icon('exclamation')
    end
    output.length>0 ? output : ""  
  end
  
  def display_information(entity, show_required=false)
    output = information(entity, show_required)
    output.length>0 ? content_tag("small", output) : ""  
  end
  
  def display_grid_next_versions(grid, entity)
    output = ""
    previous_entity_id = grid.row_select_previous_version(entity)
    next_entity_id = grid.row_select_next_version(entity)
    if previous_entity_id > 0
      output << "&nbsp;".html_safe + previous_entity_id.to_s
    end
    output << "&nbsp;".html_safe + display_information(entity) + "&nbsp;".html_safe
    if next_entity_id > 0
      output << next_entity_id.to_s + "&nbsp;".html_safe
    end
    output
  end

  def display_locale(entity)
    if entity.locale != entity.base_locale
      language = LANGUAGES.find {|lang, locale| entity.base_locale == locale}
      icon('warning') + ("&nbsp;" + link_to(language[0], 
                                           refresh_path,
                                           :locale => language[1])).html_safe
    else
      ""
    end
  end
  
  def get_language(base_locale)
    LANGUAGES.find {|lang, locale| base_locale == locale}[0]
  end
  
  def show_collection_count(collection, total, page)
    length = collection.length
    if length > 0
      if page.present? and page.to_i > 0
        from = (page.to_i - 1) * Grid::DISPLAY_ROWS_LIMIT + 1
        to = from + length - 1
      else
        from = 1
        to = length
      end
      t(length > 1 ? 'general.records' : 'general.record', 
             :count => length.to_s, 
             :total => total > 0 ? total : '?',
             :from => from,
             :to => to)
    end
  end

  def calc_next_page(page)
    page.present? ? page.to_i+1 : 2
  end

  def calc_previous_page(page)
    page.present? and page.to_i > 1 ? page.to_i-1 : 1
  end
  
  def not_first_page(page)
    page.present? and page.to_i > 1 
  end
  
  def not_all_rows(total)
    total > Grid::DISPLAY_ROWS_LIMIT
  end

  def not_last_page(page, total)
    if page.present?
      total > (page.to_i) * Grid::DISPLAY_ROWS_LIMIT_FULL
    else
      total > Grid::DISPLAY_ROWS_LIMIT_FULL
    end
  end
  
  def init_row_index(page, full=false)
    @row_index = page.present? ? ((page.to_i-1) * 
                  (full ? Grid::DISPLAY_ROWS_LIMIT_FULL : 
                          Grid::DISPLAY_ROWS_LIMIT)) : 0
  end
  
  def row_index
    @row_index += 1
    @row_index
  end

  def init_column_index
    @column_index = 0
  end
  
  def column_index
    @column_index += 1
    @column_index
  end
end