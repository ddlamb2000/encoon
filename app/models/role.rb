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
class Role < Entity
  ROOT_UUID = '62630550-0cb5-012d-bcb2-4417fe7fde95'
  
  ROLE_READ_ONLY_UUID = '50088bf0-9c87-012f-089e-4417fe7fde95'
  ROLE_READ_WRITE_UUID = '57830c40-9c87-012f-089e-4417fe7fde95'
  ROLE_READ_WRITE_ALL_UUID = '58c5b1f0-9f85-012f-1924-4417fe7fde95'
  ROLE_TOTAL_CONTROL_UUID = '69cff910-9c87-012f-089e-4417fe7fde95'
end