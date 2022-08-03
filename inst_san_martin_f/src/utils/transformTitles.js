import { LanguageTitles } from "../constant/constant";
import lodash from 'lodash';

export const transformTitles = (data) => {
  const idx = lodash.findIndex(LanguageTitles.titles, (l) => l === data);
  return idx === -1 ? data : LanguageTitles.titulos[idx]; 
};
