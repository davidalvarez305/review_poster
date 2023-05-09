import { Dictionary } from "../types/general";
import { capitalizeFirstLetter } from "./capitalizeFirstLetter";
import { convertToTitleCase } from "./convertToTitleCase";
import { getRandomInt } from "./getRandomInt";

export const re = /(\([#@]\w+:[A-Z]+)\)|(\([#@]\w+)\)/g;

function switchWords(word: string, dict: Dictionary) {
  if (!dict[word] || dict[word].length < 1) {
    return word;
  }
  return dict[word][getRandomInt(dict[word].length)];
}

export function spinContent(sentence: string, dict: Dictionary): string {
  function contentSpinner(match: string) {
    switch (true) {
      case match === "(@ProductName)":
        return (match = "Adidas Powerlift 4");
      case match === "(@ProductCategory)":
        return (match = "product category");
      default:
        // Check if the word needs to be capitalized, if not return as-is
        if (match.split(":").length === 2) {
          let word = match.split(":")[0].concat(")");
          const processed = switchWords(word, dict);
          if (match.split(":")[1] === "UU)") {
            return convertToTitleCase(processed);
          }
          if (match.split(":")[1] === "U)") {
            return capitalizeFirstLetter(processed);
          }
        }
        return switchWords(match, dict);
    }
  }
  return sentence.replaceAll(re, contentSpinner);
}
