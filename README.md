# Audiogen-CLI

A command-line tool for extracting text from PDFs, splitting it into sentences, narrating it using a [TTS API](https://github.com/erew123/alltalk_tts), and merging audio files.

## Installation

Clone the repository and build the binary:

```sh
git clone https://github.com/VinewZ/audiogen-cli
cd audiogen-cli 
make build
```

## Usage

```sh
./bin/audiogen-cli <command> [options]
```

Commands
extract (e)

Extract text from a given PDF, either saving it to the database or returning it to stdout.

```sh
./bin/audiogen-cli extract -p <path_to_pdf> [-t <book_title>]
```

Options:

```sh
    -p, --path (required): Path to the PDF file.

    -t, --title: Save extracted text and chapters to the database instead of returning it.
```

split (s)

Split text from a book into sentences.

```sh
./bin/audiogen-cli split -t <book_title> -l <language> [-c <chapter>]
```

Options:

```sh
    -t, --title (required): Title of the book.

    -l, --language (required): Language of the book.

    -c, --chapter: Split text from a specific chapter.
```

narrate (n)

Post sentences to the TTS API and retrieve the audio file.

```sh
./bin/audiogen-cli narrate -t <book_title> -l <language> [-c <chapter>]
```

Options:

```sh
    -t, --title (required): Title of the book.

    -l, --language (required): Language of the book.

    -c, --chapter: Process sentences from a specific chapter.
```

merge (m)

Merge all generated audio files into a single WAV file.

```sh
./bin/audiogen-cli merge -t <book_title> -d <delay_ms> -l <language> [-c <chapter>]
```

Options:

```sh
    -t, --title (required): Title of the book.

    -d, --delay (required): Delay between sentences in milliseconds.

    -l, --language (required): Language of the book.

    -c, --chapter: Merge audio from a specific chapter.
```
