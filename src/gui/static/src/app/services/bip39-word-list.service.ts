import { debounceTime, map } from 'rxjs/operators';
import { Subject, Observable, from } from 'rxjs';
import { Injectable } from '@angular/core';

/**
 * Allows to access the BIP39 word list.
 */
@Injectable()
export class Bip39WordListService {

  /**
   * Emits the results of the searches requested by calling the setSearchTerm function.
   */
  get searchResults(): Observable<string[]> {
    // Filter the emissions to avoid making very frecuent calculations while the user writes.
    return this.searchRequestSubject.asObservable().pipe(debounceTime(100), map(searchTerm => {
      if (searchTerm.length > 1 && this.wordList) {
        return this.wordList.filter(option => option.startsWith(searchTerm));
      } else {
        return [];
      }
    }));
  }

  private lastSearchTerm = '';
  private searchRequestSubject: Subject<string> = new Subject<string>();
  /**
   * Array with all the BIP39 words.
   */
  private wordList: string[];
  /**
   * Map with all the BIP39 words, used mainly to know if a word is on the list.
   */
  private wordMap: Map<string, boolean> = new Map<string, boolean>();

  constructor() {
    // Load the word list.
    const name = 'bip39-word-list';
    from(import(`../../assets/${name}.json`)).subscribe(result => {
      this.wordList = result.list;
      this.wordList.forEach(word => {
        this.wordMap.set(word, true);
      });

      // If a search was requested before loading the word list, the search is done again.
      this.searchRequestSubject.next(this.lastSearchTerm);
    });
  }

  /**
   * Asks the service to search for all the BIP39 words which start with the provided
   * string. The service will search for the words and the result will be emited by
   * the searchResults observable.
   * @param value String used to search for the BIP39 words.
   */
  setSearchTerm(value: string) {
    this.lastSearchTerm = value;
    if (this.wordList) {
      this.searchRequestSubject.next(value);
    }
  }

  /**
   * Checks if a string is a valid BIP39 word.
   * @param word Word to check.
   * @returns True or false, depending on whether the provided word is a valid BIP39
   * word. If the service is still loading the word list, this function returns null.
   */
  validateWord(word: string): boolean | null {
    if (this.wordList) {
      if (!this.wordMap.has(word)) {
        return false;
      }

      return true;
    } else {
      return null;
    }
  }
}
