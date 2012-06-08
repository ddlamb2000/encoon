class UpdateColumnKind < ActiveRecord::Migration
  def self.up
    Column.find(:all).each do |item|
      puts "Update column" + item.to_s
      puts "......item.kind=" + item.kind
      case item.kind
        when 'STRING' then item.kind = Column::STRING
        when 'TEXT' then item.kind = Column::TEXT
        when 'DATE' then item.kind = Column::DATE
        when 'INTEGER' then item.kind = Column::INTEGER
        when 'DECIMAL' then item.kind = Column::DECIMAL
        when 'BOOLEAN' then item.kind = Column::BOOLEAN
        when 'REFERENCE' then item.kind = Column::REFERENCE
        when 'HYPERLINK' then item.kind = Column::HYPERLINK
        when 'PASSWORD' then item.kind = Column::PASSWORD
        when 'PHOTO' then item.kind = Column::PHOTO
        when 'DOCUMENT' then item.kind = Column::DOCUMENT
      end
      puts "......=> item.kind=" + item.kind
      item.save! if item.valid?
    end
  end

  def self.down
  end
end
