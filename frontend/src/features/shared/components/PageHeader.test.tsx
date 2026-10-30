import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { PageHeader } from "./PageHeader";
import { Button } from "@mui/material";

describe("PageHeader", () => {
  it("renders semantic h1 title with variant h5 styling and subtitle", () => {
    render(
      <PageHeader
        title="Meus Veículos"
        subtitle="Gerencie sua frota com facilidade"
      />,
    );

    const heading = screen.getByRole("heading", {
      level: 1,
      name: "Meus Veículos",
    });
    expect(heading).toBeInTheDocument();
    expect(
      screen.getByText("Gerencie sua frota com facilidade"),
    ).toBeInTheDocument();
  });

  it("renders custom element in title and subtitle", () => {
    render(
      <PageHeader
        title={<span data-testid="custom-title">Título Customizado</span>}
        subtitle={
          <span data-testid="custom-subtitle">Subtítulo Customizado</span>
        }
      />,
    );

    expect(screen.getByTestId("custom-title")).toBeInTheDocument();
    expect(screen.getByTestId("custom-subtitle")).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 1 })).toBeInTheDocument();
  });

  it("renders skeleton placeholder when loading is true", () => {
    render(<PageHeader title="Veículo Teste" loading={true} />);

    expect(screen.getByTestId("page-header-skeleton")).toBeInTheDocument();
    expect(screen.queryByText("Veículo Teste")).not.toBeInTheDocument();
  });

  it("renders backButton from config object and handles click", () => {
    const onBackClick = vi.fn();
    render(
      <PageHeader
        title="Detalhes do Veículo"
        backButton={{ onClick: onBackClick, label: "Voltar para lista" }}
      />,
    );

    const backBtn = screen.getByRole("button", { name: "Voltar para lista" });
    expect(backBtn).toBeInTheDocument();
    fireEvent.click(backBtn);
    expect(onBackClick).toHaveBeenCalledTimes(1);
  });

  it("renders custom backButton element directly", () => {
    render(
      <PageHeader
        title="Editar Item"
        backButton={<button data-testid="custom-back-btn">Voltar</button>}
      />,
    );

    expect(screen.getByTestId("custom-back-btn")).toBeInTheDocument();
  });

  it("renders action slot and breadcrumbs", () => {
    render(
      <PageHeader
        title="Relatórios"
        breadcrumbs={<nav data-testid="breadcrumbs">Home &gt; Relatórios</nav>}
        action={<Button variant="contained">Novo Relatório</Button>}
      />,
    );

    expect(screen.getByTestId("breadcrumbs")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Novo Relatório" }),
    ).toBeInTheDocument();
  });
});
