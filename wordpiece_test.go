package wordpiece

import (
    "testing"
    "reflect"
)

func TestBasicTokenizeSimple(t *testing.T) {
    cases := [...]struct {
        input string
        expected []string
    }{
        {"Hello, world\tfrom    go", []string {"hello", ",", "world", "from", "go"}},
        {".., world.,", []string {".", ".", ",", "world", ".", ","}},
    }

    for _, c := range cases {
        found := BasicTokenize(c.input, true)
        if !reflect.DeepEqual(found, c.expected) {
            t.Errorf("Expected %q -> %q, found %q", c.input, c.expected, found)
        }
    }
}

func TestWordPieceTokenize(t *testing.T) {
    cases := [...]struct {
        input string
        expected []string
    }{
        {"Hello world, from go.", []string {"hello", "world", ",", "[UNK]", "go", "."}},
        {"Helloworld, from go.", []string {"hello", "[UNK]", ",", "[UNK]", "go", "."}},
        {"Hello fromgo.", []string {"hello", "[UNK]", "."}},
        {"Hello worldgo.", []string {"hello", "world", "##go", "."}},
    }

    vocab := map[string]int{
        "hello": 0,
        "world": 1,
        ",": 2,
        ".": 3,
        "go": 4,
        "##go": 5,
    }
    for _, c := range cases {
        found := WordPieceTokenize(c.input, vocab, "[UNK]", true)
        if !reflect.DeepEqual(found, c.expected) {
            t.Errorf("Expected %q -> %q, found %q", c.input, c.expected, found)
        }
    }
}

func TestBertWordPieceTokenize(t *testing.T) {
    cases := [...]struct {
        input string
        expected []string
    }{
        {"This is a pen", []string {"this", "is", "a", "pen"}},
        {"This, i-s a pen!!", []string {"this", ",", "i", "-", "s", "a", "pen", "!", "!"}},
        {"penguins are flightless birds", []string {"penguins", "are", "flight", "##less", "birds"}},
    }

    vocab := LoadVocab("/projects/tir3/users/kkurita/WeightPoisoning/bert-base-uncased-vocab.txt")
    for _, c := range cases {
        found := WordPieceTokenize(c.input, vocab, "[UNK]", true)
        if !reflect.DeepEqual(found, c.expected) {
            t.Errorf("Expected %q -> %q, found %q", c.input, c.expected, found)
        }
    }
}
