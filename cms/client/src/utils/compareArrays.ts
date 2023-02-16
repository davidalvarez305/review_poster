export const compareArrays = (arr1: Array<any>, arr2: Array<any>) => {
  if (Array.isArray(arr1) && Array.isArray(arr2)) {
    return arr1.join("") === arr2.join("");
  } else {
    return false;
  }
};
