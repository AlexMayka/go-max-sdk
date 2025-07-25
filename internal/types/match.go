package types

type MatchType int

const (
	MatchExact MatchType = iota
	MatchPrefix
	MatchSuffix
	MatchContains
	MatchRegex
	MatchCommand
	MatchCallback
	MatchBotStarted
	MatchAny
)
