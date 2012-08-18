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
Encoon::Application.routes.draw do
  root :to => 'home#index'
  resources :grids do
    resources :rows
  end  
  devise_for :users
  match '/history' => 'home#history', :as => 'history'
  match '/set' => 'home#set', :as => 'set'
  match '/unset' => 'home#unset', :as => 'unset'
  match '/import' => 'home#import', :as => 'import'
  match '/refresh' => 'home#refresh', :as => 'refresh'
  match ':controller/:action/:id'
  match ':controller/:action/:id.:format'
end