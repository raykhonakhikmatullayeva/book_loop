# book_loop

[![Lint](https://github.com/raykhonakhikmatullayeva/book_loop/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/raykhonakhikmatullayeva/book_loop/actions/workflows/lint.yml)
[![Nilaway](https://github.com/raykhonakhikmatullayeva/book_loop/actions/workflows/nilaway.yml/badge.svg?branch=main)](https://github.com/raykhonakhikmatullayeva/book_loop/actions/workflows/nilaway.yml)
[![Coverage](https://github.com/raykhonakhikmatullayeva/book_loop/blob/badges/coverage.svg)](https://github.com/raykhonakhikmatullayeva/book_loop/blob/badges/coverage.svg)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)

## Описание

**BookLoop** — онлайн-библиотека. Читатель ищет книги в каталоге и берёт их на дом; 
библиотекарь ведёт каталог, отмечает возвраты и потери, следит за сроками. 
Ключевое ограничение продукта — число экземпляров конечно: когда книгу разобрали, система обязана вести себя корректно — 
ставить в очередь и следить за просрочками.