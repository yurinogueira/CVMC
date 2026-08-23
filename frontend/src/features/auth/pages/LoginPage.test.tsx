import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { LoginPage } from "./LoginPage";

describe("LoginPage", () => {
  it("renders login form fields, labels, and CTA buttons", () => {
    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>,
    );

    expect(screen.getByText("Bem-vindo de volta")).toBeInTheDocument();
    expect(screen.getByLabelText(/^E-mail/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Senha/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Entrar no CVMC/i }),
    ).toBeInTheDocument();
    expect(screen.getByText(/Criar conta gratuita/i)).toBeInTheDocument();
  });
});
