package wordpiece

import (
    "bufio"
    "log"
    "strings"
    "io"
    "os"
    "unicode"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

const MaxTokenLen = 100


func LoadVocab(path string) map[string]int {
    file, err := os.Open(path) 
    if err != nil { log.Fatal(err) }
    defer file.Close()

    vocab := make(map[string]int)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        text := scanner.Text()
        vocab[text] = len(vocab)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return vocab
}

func whitespace_tokenize(s string) []string {
    return strings.Fields(strings.TrimSpace(s))
}

func clean(token string, do_lower_case bool) string {
    if do_lower_case {
        token = strings.ToLower(token)
    }
    var sb strings.Builder
    for _, c := range token {
        if unicode.IsSpace(c) {
            sb.WriteRune(' ')
        } else if unicode.Is(unicode.Mn, c) {
            // skip accented characters
        } else if (c == 0 || c == 0xfffd || unicode.IsControl(c)) {
            // skip control characters 
        } else {
            sb.WriteRune(c)
        }
    }
    token = sb.String()
    return token
}

/**
* Reads in data while normalizing unicode
*/
func NormalizedReader(r io.Reader) io.Reader {
    t := transform.Chain(norm.NFD, norm.NFC)
    return transform.NewReader(r, t)
}

func _split_on_punc(text string) []string {
    var words []string
    var sb strings.Builder

    start_new_word := true
    for _, c := range text {
        if unicode.IsPunct(c) {
            if !start_new_word {
                words = append(words, sb.String())
                sb.Reset()
            }
            words = append(words, string(c))
            start_new_word = true
        } else {
            sb.WriteRune(c)
            start_new_word = false
        }
    }
    if !start_new_word {
        words = append(words, sb.String())
        sb.Reset()
    }
    return words
}

func BasicTokenize(text string, do_lower_case bool) []string {
    orig_tokens := whitespace_tokenize(text)
    var output_tokens []string
    for _, token := range orig_tokens {
        output_tokens = append(output_tokens, _split_on_punc(clean(token, do_lower_case))...)
    }
    return output_tokens
}

func in_vocab(token string, vocab map[string] int) bool {
    _, exists := vocab[token]
    return exists
}

func subword_tokenize(token string, vocab map[string]int, unk_token string) []string {
    var sub_tokens []string
    if len(token) > MaxTokenLen {
        sub_tokens = append(sub_tokens, unk_token)
        return sub_tokens
    }
    start := 0
    for start < len(token) {
        end := len(token)
        substr_is_valid := false
        var cur_substr string
        for start < end { 
            cur_substr = string(token[start:end])
            if start > 0 {
                cur_substr = "##" + cur_substr
            }
            if in_vocab(cur_substr, vocab) {
                substr_is_valid = true
                break
            }
            end -= 1
        }
        if !substr_is_valid {
            sub_tokens = append(sub_tokens, unk_token)
            break
        }
        sub_tokens = append(sub_tokens, cur_substr)
        start = end
    }
    return sub_tokens
}

func WordPieceTokenize(text string, vocab map[string]int, unk_token string, do_lower_case bool) []string {
    var output_tokens []string
    orig_tokens := BasicTokenize(text, do_lower_case)
    for _, token := range orig_tokens {
        output_tokens = append(output_tokens, subword_tokenize(token, vocab, unk_token)...)
    }
    return output_tokens
}
