module github.com/cmburn/perlutils

go 1.19

// I'm aware this is an outdated version, however in v7.14.0 they introduced a
// "product check", which breaks CPAN.
require github.com/elastic/go-elasticsearch/v7 v7.13.1

require github.com/relvacode/iso8601 v1.1.0
