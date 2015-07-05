#!/usr/bin/env ruby

require 'tmpdir'
dir = Dir.mktmpdir
Dir.chdir(dir)

# Create a repo with some history
`got init`
File.write('file.txt', 'First revision.')
`got commit -m "First commit" -a "Testing Script"`
File.write('file.txt', 'Second revision.')
`got commit -m "Second commit" -a "Testing Script"`
File.write('file.txt', 'Third revision.')

# Revert the file
`got revert file.txt`

if IO.read("file.txt") != 'Second revision.'
  exit 1
end
