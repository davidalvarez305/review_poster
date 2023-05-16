import { CreateWordInput, Word, WordFormInput } from "../types/general";
import createTagFactory from "./createTagFactory";
import { getId } from "./getId";

export default function createWordsFactory(values: WordFormInput, user_id: number, words: Word[]): CreateWordInput {
    return { ...values, id: getId(values.word, words, "name"), tag: createTagFactory(values.word), word: values.word , user_id, synonyms: values.synonyms.split("\n")};
}