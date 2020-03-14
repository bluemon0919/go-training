package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf8"
)

type RuneScanner struct {
	r   io.Reader
	buf [16]byte
	c   rune
	err error
}

func NewRuneScanner(r io.Reader) *RuneScanner {
	return &RuneScanner{r: r}
}

func (s *RuneScanner) Scan() bool {
	if s.err != nil {
		return false
	}

	// buf[16], n:13
	// buf[:n] // データが入ってる部分
	n, err := s.r.Read(s.buf[:])
	if err != nil && err != io.EOF {
		s.err = err
		return false
	}
	if err == io.EOF {
		return false
	}

	fmt.Println(n, "bytes read")

	// buf[:n]は[]byte型なのでバイト列
	// Goの文字列はUTF-8でエンコードされている
	// buf[:n]の中から先頭のルーンを取り出す
	// sizeは先頭のルーンに何バイト使うかが返ってくる
	c, size := utf8.DecodeRune(s.buf[:n]) // 足りなかった場合の処理が抜けてる
	if c == utf8.RuneError {
		s.err = errors.New("RuneError")
		return false
	}
	s.c = c

	s.r = io.MultiReader(bytes.NewReader(s.buf[size:n]), s.r)
	return true
}

func (s *RuneScanner) Rune() rune { // ガワを先に決めておく
	return s.c
}

func (s *RuneScanner) Err() error {
	return s.err
}

func main() {
	s := NewRuneScanner(strings.NewReader("Hello, 世界"))
	for s.Scan() { // bool
		r := s.Rune() // rune
		fmt.Printf("%c\n", r)
	}
	if err := s.Err(); err != nil { // error
		log.Fatal(err)
	}
}
