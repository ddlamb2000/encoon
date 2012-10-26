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
class Cache < Entity
  # Internal cache used for storing loaded grid definitions.
  @@grid_cache = []

  # Returns a flat string with the content of a filter.
  def self.flat_filters(filters)
    return "" if filters.nil?
    flat = "" 
    for filter in filters
      flat << filter[:column_uuid] + ":" + filter[:row_uuid]  
    end
    flat
  end

  # Returns loaded grid information from the internal grid cache.
  def self.get_cached_grid(uuid, workspace_uuid=nil, uri=nil, filters=nil)
    log_debug "Cache#get_cached_grid(#{uuid}, #{workspace_uuid}, #{uri}, #{flat_filters(filters)})"
    cached = @@grid_cache.find {|value| value[:user_uuid] == Entity.session_user_uuid and
                                        value[:asofdate] == Entity.session_as_of_date and
                                        value[:locale] == Entity.session_locale and
                                        ((uuid.present? and
                                          value[:uuid] == uuid) or 
                                         (uri.present? and
                                          workspace_uuid.present? and
                                          value[:workspace_uuid] == workspace_uuid and
                                          value[:uri] == uri)) and
                                        value[:filters] == flat_filters(filters)}
    if cached.present?
      log_debug "Cache#get_cached_grid found"
      log_cache_grid cached
      return cached[:grid]
    end
    log_debug "Cache#get_cached_grid not found"
    nil
  end

  # Pushes loaded grid information into the internal grid cache.
  def self.grid_cache_push(grid, filters)
    log_debug "Cache#grid_cache_push(#{grid.to_s}, #{flat_filters(filters)})"
    cached = get_cached_grid(grid.uuid, filters)
    if cached.nil?
      @@grid_cache << {:user_uuid => Entity.session_user_uuid,
                       :asofdate => Entity.session_as_of_date,
                       :locale => Entity.session_locale,
                       :uuid => grid.uuid,
                       :uri => grid.uri,
                       :workspace_uuid => grid.workspace_uuid,
                       :filters => flat_filters(filters),
                       :grid => grid}
      log_debug "Cache#grid_cache_push pushed"
      for cache in @@grid_cache
        log_cache_grid cache
      end
    end
  end

  # Logs information about the cached grid.
  def self.log_cache_grid(cache)
    log_debug "Cache#log_cache_grid " +
              "uuid=#{cache[:uuid]}, " +
              "grid=#{cache[:grid]}, " +
              "asofdate=#{cache[:asofdate]}, " +
              "locale=#{cache[:locale]}, " +
              "filters=#{cache[:filters]}"
  end
end