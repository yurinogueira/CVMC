import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { PrivacyPolicyPage } from "./PrivacyPolicyPage";

describe("PrivacyPolicyPage", () => {
  it("renders public privacy policy and account deletion instructions without authentication", () => {
    render(
      <BrowserRouter>
        <PrivacyPolicyPage />
      </BrowserRouter>,
    );

    expect(
      screen.getByRole("heading", { name: /Política de Privacidade/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", {
        name: /4\. Exclusão de Conta e Dados do Usuário/i,
      }),
    ).toBeInTheDocument();
    expect(screen.getByText(/com\.cvmc\.app/i)).toBeInTheDocument();
  });
});
