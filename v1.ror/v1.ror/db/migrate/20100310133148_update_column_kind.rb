class UpdateColumnKind < ActiveRecord::Migration
  def self.up
    Column.find(:all).each do |item|
      puts "Update column" + item.to_s
      puts "......item.kind=" + item.kind
      case item.kind
        when 'STRING' then item.kind = COLUMN_TYPE_STRING
        when 'TEXT' then item.kind = COLUMN_TYPE_TEXT
        when 'DATE' then item.kind = COLUMN_TYPE_DATE
        when 'INTEGER' then item.kind = COLUMN_TYPE_INTEGER
        when 'DECIMAL' then item.kind = COLUMN_TYPE_DECIMAL
        when 'BOOLEAN' then item.kind = COLUMN_TYPE_BOOLEAN
        when 'REFERENCE' then item.kind = COLUMN_TYPE_REFERENCE
        when 'HYPERLINK' then item.kind = COLUMN_TYPE_HYPERLINK
        when 'PASSWORD' then item.kind = COLUMN_TYPE_PASSWORD
        when 'PHOTO' then item.kind = COLUMN_TYPE_PHOTO
        when 'DOCUMENT' then item.kind = COLUMN_TYPE_DOCUMENT
      end
      puts "......=> item.kind=" + item.kind
      item.save! if item.valid?
    end
  end

  def self.down
  end
end
