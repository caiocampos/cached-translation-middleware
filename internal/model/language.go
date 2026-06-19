package model

type Language string

const (
	LanguagePortuguese Language = "pt"
	LanguageEnglish    Language = "en"
	LanguageSpanish    Language = "es"
	LanguageItalian    Language = "it"
	LanguageFrench     Language = "fr"
	LanguageGerman     Language = "de"
)

const SourceLanguage = LanguagePortuguese

var TargetLanguages = [...]Language{
	LanguageEnglish,
	LanguageSpanish,
	LanguageItalian,
	LanguageFrench,
	LanguageGerman,
}
