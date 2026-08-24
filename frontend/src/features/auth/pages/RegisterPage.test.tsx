import { describe, expect, it } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { RegisterPage } from "./RegisterPage";

describe("RegisterPage", () => {
  it("renders register form fields and actions", () => {
    render(
      <BrowserRouter>
        <RegisterPage />
      </BrowserRouter>,
    );

    expect(screen.getByText("Criar uma conta")).toBeInTheDocument();
    expect(screen.getByLabelText(/Nome completo/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^E-mail/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Senha/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Confirmar Senha/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Concluir Cadastro/i }),
    ).toBeInTheDocument();
  });

  it("shows error message when password has less than 8 characters", async () => {
    render(
      <BrowserRouter>
        <RegisterPage />
      </BrowserRouter>,
    );

    fireEvent.change(screen.getByLabelText(/Nome completo/i), {
      target: { value: "Test User" },
    });
    fireEvent.change(screen.getByLabelText(/^E-mail/i), {
      target: { value: "test@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/^Senha/i), {
      target: { value: "short" },
    });
    fireEvent.change(screen.getByLabelText(/Confirmar Senha/i), {
      target: { value: "short" },
    });

    fireEvent.click(screen.getByRole("button", { name: /Concluir Cadastro/i }));

    expect(
      await screen.findByText("A senha deve conter no mínimo 8 caracteres."),
    ).toBeInTheDocument();
  });

  it("shows error message when passwords do not match", async () => {
    render(
      <BrowserRouter>
        <RegisterPage />
      </BrowserRouter>,
    );

    fireEvent.change(screen.getByLabelText(/Nome completo/i), {
      target: { value: "Test User" },
    });
    fireEvent.change(screen.getByLabelText(/^E-mail/i), {
      target: { value: "test@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/^Senha/i), {
      target: { value: "validpassword123" },
    });
    fireEvent.change(screen.getByLabelText(/Confirmar Senha/i), {
      target: { value: "differentpassword123" },
    });

    fireEvent.click(screen.getByRole("button", { name: /Concluir Cadastro/i }));

    expect(
      await screen.findByText(
        "As senhas não conferem. Verifique e tente novamente.",
      ),
    ).toBeInTheDocument();
  });
});
