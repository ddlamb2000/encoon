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
  # Returns the title of the page to be displayed in the browser.
  def current_page_title
    @page_title.present? ? (@page_title + " | " + current_application_title) : current_application_title 
  end

  # Returns the name of the application and the environment running it.
  def current_application_title
    APPLICATION_TITLE + ('production' != Rails.env ? (" (" + Rails.env + ")") : "")
  end

  # Displays an icon using its name
  def icon(name, title=nil)
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

  # Displays the filter in plain text.
  def display_filters(grid, filters, search=nil)
    output = ""
    if filters.present?
      filters.collect do |filter|
        column_uuid = filter[:column_uuid]
        row_name = filter[:row_name]
        column = grid.column_select_entity_by_uuid(column_uuid)
        if column.present? and column.kind == COLUMN_TYPE_REFERENCE
          output << ", " if output != ""
          output << t('filters.equals_to', :column => column.name, :value => row_name)
        end
      end
    end
    if search.present?
      output << ", " if output != ""
      output << t('filters.search', :value => search)
    end
    output
  end

  # Returns a string made of identifiers based on filters, 
  # used to identify a grid and associated filters in a page.
  def filters_uuid(entity, filters)
    output = entity.uuid + "-"
    filters.collect{|filter| output << filter[:column_uuid]} if filters.present?
    output
  end

  # Displays the history of navigation in plain text.
  def display_history(history)
    output = ""
    if history.present?
      history.reverse_each do |link|
        hyperlink = content_tag("a", link[:page_title].html_safe, :href => link[:url])
        hyperlink = link[:page_title] if link[:url] == request.url
        output << content_tag("li",
                              hyperlink + " " + 
                              content_tag(
                                "span",
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

  # Displays the distance between the date and the current system date.
  def display_distance_date(date)
    if date != Entity.begin_of_time and date != Entity.end_of_time
      now = Time.now
      today = Date::civil(now.year, now.month, now.day)
      t(date > today ? 'general.ahead' : 'general.ago', :time => time_ago_in_words(date))
    end
  end
  
  def display_new(date=nil)
    (date.nil? or ((Time.now-date) < 1.day)) ? 
      content_tag("div", t('general.new'), :class => "new") :
      ""
  end

  # Displays timing about the creation or last update of the entity.
  def display_updated_date(entity)
    t(entity.revision == 1 ? 'general.created' : 'general.updated', 
      :time => time_ago_in_words(entity.updated_at, :include_seconds => true)) +
    display_new(entity.updated_at)
  end

  # Displays information about the creation of the entity,
  # including timing and user who created the data.
  def display_created_time_by(entity, who, who_uuid)
    t('general.time_by',
      :time => time_ago_in_words(entity.created_at, :include_seconds => true),
      :by => (who.present? and who_uuid.present?) ?
               link_to_unless_current(who,
                 show_path(:workspace => SECURITY_WORKSPACE_URI,
                           :grid => USER_URI,
                           :row => who_uuid)) :
               t('general.unknown')).html_safe
  end

  # Displays information about the update of the entity,
  # including timing and user who last updated the data.
  def display_updated_time_by(entity, who, who_uuid)
    (t('general.time_by',
      :time => time_ago_in_words(entity.updated_at, :include_seconds => true),
      :by => (who.present? and who_uuid.present?) ?
               link_to_unless_current(who,
                 show_path(:workspace => SECURITY_WORKSPACE_URI,
                           :grid => USER_URI,
                           :row => who_uuid)) :
               t('general.unknown')) +
    display_new(entity.updated_at)).html_safe
  end
  
  # Displays information about the creation or the last update of the entity,
  # including timing and user who made the creation or last update.
  def display_updated_date_by(entity)
    if entity.revision.nil? or entity.revision == 1
      (t('field.created') + " " + display_created_time_by(entity, entity.who_created, entity.create_user_uuid)).html_safe
    else
      (t('field.updated') + " " + display_updated_time_by(entity, entity.who_updated, entity.update_user_uuid)).html_safe
    end
  end

  def warning_current_date(as_of_date)
    now = Time.now
    today = Date::civil(now.year, now.month, now.day)
    as_of_date != today ? "warning" : ""
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
    output.length > 0 ? output : ""  
  end

  def display_information(entity, show_required=false)
    output = information(entity, show_required)
    output.length > 0 ? content_tag("span", output, :class => "asofdate").html_safe : ""
  end

  def display_grid_next_versions(grid, entity)
    output = ""
    previous_entity_id = grid.row_select_previous_version(entity)
    next_entity_id = grid.row_select_next_version(entity)
    if previous_entity_id > 0
      output << " " + previous_entity_id.to_s
    end
    output << " " + display_information(entity) + " "
    if next_entity_id > 0
      output << next_entity_id.to_s + " "
    end
    output.html_safe
  end

  # Displays the language in which the row was captured if this doesn't 
  # correspond to the language being displayed. 
  def display_locale(row)
    if row.locale != row.base_locale
      language = LANGUAGES.find{|lang, locale| row.base_locale == locale}
      url = request.url.gsub(/[?&]locale=(..)/, "")
      url = url + (url["?"].nil? ? "?" : "&") + "locale=" + language[1]
      content_tag("span", link_to(language[0], url), :class => 'warning-alert')
    end
  end

  # Returns the name of the displayed language for the given locale code.
  def get_language(base_locale)
    LANGUAGES.find{|lang, locale| base_locale == locale}[0]
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
end