import { useEffect } from "react";

const BASE_TITLE = "CVMC - Como Vai Meu Carro";

export function useDocumentTitle(title?: string): void {
  useEffect(() => {
    if (title) {
      document.title = `${title} | ${BASE_TITLE}`;
    } else {
      document.title = `${BASE_TITLE} | Gestão Inteligente de Veículos e Manutenções`;
    }
  }, [title]);
}
