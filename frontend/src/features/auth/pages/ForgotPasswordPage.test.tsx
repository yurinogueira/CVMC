import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { ForgotPasswordPage } from "./ForgotPasswordPage";

describe("ForgotPasswordPage", () => {
  it("renders forgot password form and CTA button", () => {
    render(
      <BrowserRouter>
        <ForgotPasswordPage />
      </BrowserRouter>,
    );

    expect(screen.getByText("Recuperar Senha")).toBeInTheDocument();
    expect(screen.getByLabelText(/^E-mail cadastrado/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Enviar Link de Recuperação/i }),
    ).toBeInTheDocument();
    expect(screen.getByText(/Voltar para o Login/i)).toBeInTheDocument();
  });
});
