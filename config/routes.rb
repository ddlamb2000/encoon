# encoding: utf-8
Encoon::Application.routes.draw do
  root :to => 'home#index'
  resources :grids do
    resources :rows
  end  
  resources :rows
  match '/login' => 'home#login'
  match '/logout' => 'home#logout'
  match '/about' => 'home#about'
  match '/hide_about' => 'home#hide_about'
  match '/history' => 'home#history'
  match '/hide_history' => 'home#hide_history'
  match '/register' => 'home#register'
  match '/export_system' => 'home#export_system'
  match '/import' => 'home#import'
  match '/refresh' => 'home#refresh'
  match ':controller/:action/:id'
  match ':controller/:action/:id.:format'
end