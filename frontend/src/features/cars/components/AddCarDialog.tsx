import React, { useState, useEffect } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Stack,
  Alert,
  CircularProgress,
  Grid,
  Autocomplete,
  ToggleButtonGroup,
  ToggleButton,
  Typography,
  Box,
  Paper,
  Chip,
} from "@mui/material";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";
import TwoWheelerRoundedIcon from "@mui/icons-material/TwoWheelerRounded";
import LocalShippingRoundedIcon from "@mui/icons-material/LocalShippingRounded";
import CheckCircleOutlineRoundedIcon from "@mui/icons-material/CheckCircleOutlineRounded";
import MonetizationOnRoundedIcon from "@mui/icons-material/MonetizationOnRounded";
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import TagRoundedIcon from "@mui/icons-material/TagRounded";
import { ImageUploadField } from "./ImageUploadField";

import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../../auth/state/auth.store";
import { carService } from "../services/car.service";
import { fipeService } from "../services/fipe.service";
import { Car } from "../types/car.types";
import {
  FipeBrand,
  FipeModel,
  FipeVehicleDetail,
  FipeYear,
  VehicleType,
} from "../types/fipe.types";

interface AddCarDialogProps {
  open: boolean;
  onClose: () => void;
  onCarCreated: (car: Car) => void;
}

export function AddCarDialog({
  open,
  onClose,
  onCarCreated,
}: AddCarDialogProps) {
  const currentYear = new Date().getFullYear();
  const { user } = useAuthStore();
  const navigate = useNavigate();

  // Vehicle type
  const [vehicleType, setVehicleType] = useState<VehicleType>("cars");

  // Cascade selections
  const [brands, setBrands] = useState<FipeBrand[]>([]);
  const [selectedBrand, setSelectedBrand] = useState<FipeBrand | null>(null);

  const [models, setModels] = useState<FipeModel[]>([]);
  const [selectedModel, setSelectedModel] = useState<FipeModel | null>(null);

  const [years, setYears] = useState<FipeYear[]>([]);
  const [selectedYear, setSelectedYear] = useState<FipeYear | null>(null);

  const [fipeDetail, setFipeDetail] = useState<FipeVehicleDetail | null>(null);

  // Loading states
  const [loadingBrands, setLoadingBrands] = useState(false);
  const [loadingModels, setLoadingModels] = useState(false);
  const [loadingYears, setLoadingYears] = useState(false);
  const [loadingDetail, setLoadingDetail] = useState(false);

  // Form fields
  const [name, setName] = useState("");
  const [manufacturer, setManufacturer] = useState("");
  const [model, setModel] = useState("");
  const [yearManufacture, setYearManufacture] = useState<number>(currentYear);
  const [yearModel, setYearModel] = useState<number>(currentYear);
  const [lastMileage, setLastMileage] = useState<number>(0);
  const [imageUrl, setImageUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const resetForm = () => {
    setVehicleType("cars");
    setSelectedBrand(null);
    setSelectedModel(null);
    setSelectedYear(null);
    setFipeDetail(null);
    setBrands([]);
    setModels([]);
    setYears([]);
    setName("");
    setManufacturer("");
    setModel("");
    setYearManufacture(currentYear);
    setYearModel(currentYear);
    setLastMileage(0);
    setImageUrl("");
    setErrorMsg(null);
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const handleVehicleTypeChange = (newType: VehicleType) => {
    setVehicleType(newType);
    setSelectedBrand(null);
    setSelectedModel(null);
    setSelectedYear(null);
    setFipeDetail(null);
    setModels([]);
    setYears([]);
  };

  const handleBrandChange = (newBrand: FipeBrand | null) => {
    setSelectedBrand(newBrand);
    setSelectedModel(null);
    setSelectedYear(null);
    setFipeDetail(null);
    setModels([]);
    setYears([]);
  };

  const handleModelChange = (newModel: FipeModel | null) => {
    setSelectedModel(newModel);
    setSelectedYear(null);
    setFipeDetail(null);
    setYears([]);
  };

  const handleYearChange = (newYear: FipeYear | null) => {
    setSelectedYear(newYear);
    if (!newYear) {
      setFipeDetail(null);
    }
  };

  // Load brands when dialog opens or vehicleType changes
  useEffect(() => {
    if (!open) return;

    let isMounted = true;
    const fetchBrands = async () => {
      try {
        setLoadingBrands(true);
        const data = await fipeService.getBrands(vehicleType);
        if (isMounted) {
          setBrands(data);
        }
      } catch {
        if (isMounted) {
          setErrorMsg("Erro ao carregar lista de marcas da Fipe.");
        }
      } finally {
        if (isMounted) {
          setLoadingBrands(false);
        }
      }
    };

    fetchBrands();
    return () => {
      isMounted = false;
    };
  }, [open, vehicleType]);

  // Load models when brand is selected
  useEffect(() => {
    if (!selectedBrand) return;

    let isMounted = true;
    const fetchModels = async () => {
      try {
        setLoadingModels(true);
        const data = await fipeService.getModels(
          vehicleType,
          selectedBrand.code,
        );
        if (isMounted) {
          setModels(data);
        }
      } catch {
        if (isMounted) {
          setErrorMsg("Erro ao carregar modelos da marca selecionada.");
        }
      } finally {
        if (isMounted) {
          setLoadingModels(false);
        }
      }
    };

    fetchModels();
    return () => {
      isMounted = false;
    };
  }, [selectedBrand, vehicleType]);

  // Load years when model is selected
  useEffect(() => {
    if (!selectedBrand || !selectedModel) return;

    let isMounted = true;
    const fetchYears = async () => {
      try {
        setLoadingYears(true);
        const data = await fipeService.getYears(
          vehicleType,
          selectedBrand.code,
          selectedModel.code,
        );
        if (isMounted) {
          setYears(data);
        }
      } catch {
        if (isMounted) {
          setErrorMsg("Erro ao carregar anos do modelo selecionado.");
        }
      } finally {
        if (isMounted) {
          setLoadingYears(false);
        }
      }
    };

    fetchYears();
    return () => {
      isMounted = false;
    };
  }, [selectedBrand, selectedModel, vehicleType]);

  // Load vehicle details when year is selected
  useEffect(() => {
    if (!selectedBrand || !selectedModel || !selectedYear) return;

    let isMounted = true;
    const fetchDetail = async () => {
      try {
        setLoadingDetail(true);
        const detail = await fipeService.getVehicleDetail(
          vehicleType,
          selectedBrand.code,
          selectedModel.code,
          selectedYear.code,
        );
        if (isMounted) {
          setFipeDetail(detail);
          setManufacturer(detail.brand || selectedBrand.name);
          setModel(detail.model || selectedModel.name);
          const yearVal = detail.modelYear || parseInt(selectedYear.code, 10);
          setYearModel(yearVal || currentYear);
          setYearManufacture(yearVal || currentYear);
          setName((prev) => (prev ? prev : detail.model || selectedModel.name));
        }
      } catch {
        if (isMounted) {
          setErrorMsg("Erro ao carregar detalhes e valor Fipe do veículo.");
        }
      } finally {
        if (isMounted) {
          setLoadingDetail(false);
        }
      }
    };

    fetchDetail();
    return () => {
      isMounted = false;
    };
  }, [selectedBrand, selectedModel, selectedYear, vehicleType, currentYear]);

  const handleSubmit = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setErrorMsg(null);

    if (!name.trim() || !manufacturer.trim() || !model.trim()) {
      setErrorMsg(
        "Selecione o veículo na Tabela Fipe e preencha o apelido identificador.",
      );
      return;
    }

    if (Number(lastMileage) < 0) {
      setErrorMsg("A quilometragem não pode ser negativa.");
      return;
    }

    if (
      Number(yearManufacture) < 1900 ||
      Number(yearManufacture) > currentYear + 1 ||
      Number(yearModel) < 1900 ||
      Number(yearModel) > currentYear + 2
    ) {
      setErrorMsg("Informe anos de fabricação e modelo válidos.");
      return;
    }

    try {
      setLoading(true);
      const newCar = await carService.create({
        name: name.trim(),
        manufacturer: manufacturer.trim(),
        model: model.trim(),
        yearManufacture: Number(yearManufacture),
        yearModel: Number(yearModel),
        lastMileage: Number(lastMileage) || 0,
        vehicleType,
        imageUrl: imageUrl.trim() || undefined,
        fipeCode: fipeDetail?.codeFipe,
        fipePrice: fipeDetail?.price,
        fuel: fipeDetail?.fuel,
      });

      onCarCreated(newCar);
      handleClose();
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível cadastrar o veículo. Tente novamente.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      fullWidth
      maxWidth="md"
      slotProps={{
        paper: {
          elevation: 0,
          sx: {
            borderRadius: 3,
            border: "1px solid #E2E8F0",
            p: 1,
          },
        },
      }}
    >
      <DialogTitle sx={{ fontWeight: 700, pb: 1 }}>
        Cadastrar Novo Veículo (Tabela FIPE)
      </DialogTitle>

      <DialogContent dividers sx={{ borderColor: "#F1F5F9" }}>
        {user?.emailVerified === false && (
          <Alert
            severity="warning"
            action={
              <Button
                color="inherit"
                size="small"
                onClick={() => {
                  handleClose();
                  navigate("/profile");
                }}
              >
                Ver Perfil
              </Button>
            }
            sx={{ mb: 2.5, borderRadius: 1.5 }}
          >
            Para cadastrar veículos, é obrigatório validar seu endereço de
            e-mail. Acesse seu perfil para confirmar.
          </Alert>
        )}

        {errorMsg && (
          <Alert
            severity="error"
            onClose={() => setErrorMsg(null)}
            sx={{ mb: 2.5, borderRadius: 1.5 }}
          >
            {errorMsg}
          </Alert>
        )}

        <Stack spacing={3} sx={{ mt: 1 }}>
          {/* Tipo de Veículo */}
          <Box>
            <Typography
              variant="caption"
              sx={{
                fontWeight: 600,
                color: "text.secondary",
                textTransform: "uppercase",
                letterSpacing: 0.5,
                display: "block",
                mb: 1,
              }}
            >
              1. Tipo de Veículo
            </Typography>
            <ToggleButtonGroup
              value={vehicleType}
              exclusive
              onChange={(_, nextType) => {
                if (nextType !== null) {
                  handleVehicleTypeChange(nextType);
                }
              }}
              aria-label="Tipo de veículo"
              fullWidth
              size="small"
            >
              <ToggleButton value="cars" aria-label="Carros" sx={{ py: 1 }}>
                <DirectionsCarRoundedIcon sx={{ mr: 1, fontSize: 20 }} />
                Carros
              </ToggleButton>
              <ToggleButton
                value="motorcycles"
                aria-label="Motos"
                sx={{ py: 1 }}
              >
                <TwoWheelerRoundedIcon sx={{ mr: 1, fontSize: 20 }} />
                Motos
              </ToggleButton>
              <ToggleButton
                value="trucks"
                aria-label="Caminhões"
                sx={{ py: 1 }}
              >
                <LocalShippingRoundedIcon sx={{ mr: 1, fontSize: 20 }} />
                Caminhões
              </ToggleButton>
            </ToggleButtonGroup>
          </Box>

          {/* Seleção em Cascata Fipe */}
          <Box>
            <Typography
              variant="caption"
              sx={{
                fontWeight: 600,
                color: "text.secondary",
                textTransform: "uppercase",
                letterSpacing: 0.5,
                display: "block",
                mb: 1.5,
              }}
            >
              2. Consulta FIPE em Cascata
            </Typography>

            <Grid container spacing={2}>
              {/* Marca */}
              <Grid size={{ xs: 12, md: 4 }}>
                <Autocomplete
                  options={brands}
                  getOptionLabel={(option) => option.name}
                  getOptionKey={(option) => option.code}
                  isOptionEqualToValue={(option, value) =>
                    option.code === value.code
                  }
                  value={selectedBrand}
                  onChange={(_, newValue) => handleBrandChange(newValue)}
                  loading={loadingBrands}
                  disabled={loadingBrands || loading}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Marca"
                      placeholder="Busque a marca..."
                      required
                      slotProps={{
                        ...params.slotProps,
                        input: {
                          ...params.slotProps.input,
                          endAdornment: (
                            <>
                              {loadingBrands ? (
                                <CircularProgress color="inherit" size={18} />
                              ) : null}
                              {params.slotProps.input.endAdornment}
                            </>
                          ),
                        },
                      }}
                    />
                  )}
                />
              </Grid>

              {/* Modelo */}
              <Grid size={{ xs: 12, md: 4 }}>
                <Autocomplete
                  options={models}
                  getOptionLabel={(option) => option.name}
                  getOptionKey={(option) => option.code}
                  isOptionEqualToValue={(option, value) =>
                    option.code === value.code
                  }
                  value={selectedModel}
                  onChange={(_, newValue) => handleModelChange(newValue)}
                  loading={loadingModels}
                  disabled={!selectedBrand || loadingModels || loading}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Modelo"
                      placeholder={
                        selectedBrand
                          ? "Selecione o modelo..."
                          : "Selecione a marca primeiro"
                      }
                      required
                      slotProps={{
                        ...params.slotProps,
                        input: {
                          ...params.slotProps.input,
                          endAdornment: (
                            <>
                              {loadingModels ? (
                                <CircularProgress color="inherit" size={18} />
                              ) : null}
                              {params.slotProps.input.endAdornment}
                            </>
                          ),
                        },
                      }}
                    />
                  )}
                />
              </Grid>

              {/* Ano / Combustível / Versão */}
              <Grid size={{ xs: 12, md: 4 }}>
                <Autocomplete
                  options={years}
                  getOptionLabel={(option) => option.name}
                  getOptionKey={(option) => option.code}
                  isOptionEqualToValue={(option, value) =>
                    option.code === value.code
                  }
                  value={selectedYear}
                  onChange={(_, newValue) => handleYearChange(newValue)}
                  loading={loadingYears}
                  disabled={!selectedModel || loadingYears || loading}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Modelo / Versão"
                      placeholder={
                        selectedModel
                          ? "Selecione o modelo / versão..."
                          : "Selecione o modelo primeiro"
                      }
                      required
                      slotProps={{
                        ...params.slotProps,
                        input: {
                          ...params.slotProps.input,
                          endAdornment: (
                            <>
                              {loadingYears || loadingDetail ? (
                                <CircularProgress color="inherit" size={18} />
                              ) : null}
                              {params.slotProps.input.endAdornment}
                            </>
                          ),
                        },
                      }}
                    />
                  )}
                />
              </Grid>
            </Grid>
          </Box>

          {/* Card de Dados Fipe Autopreenchidos */}
          {fipeDetail && (
            <Paper
              elevation={0}
              sx={{
                p: 2,
                borderRadius: 2,
                bgcolor: "action.hover",
                border: "1px solid",
                borderColor: "divider",
              }}
            >
              <Stack
                direction="row"
                spacing={1}
                sx={{ mb: 1.5, alignItems: "center" }}
              >
                <CheckCircleOutlineRoundedIcon
                  color="success"
                  fontSize="small"
                />
                <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                  Dados FIPE Identificados
                </Typography>
                <Chip
                  label={fipeDetail.referenceMonth}
                  size="small"
                  variant="outlined"
                  sx={{ fontSize: "0.75rem", ml: "auto" }}
                />
              </Stack>

              <Grid container spacing={2}>
                <Grid size={{ xs: 12, sm: 4 }}>
                  <Stack
                    direction="row"
                    spacing={1}
                    sx={{ alignItems: "center" }}
                  >
                    <TagRoundedIcon
                      fontSize="small"
                      sx={{ color: "text.secondary" }}
                    />
                    <Box>
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ display: "block" }}
                      >
                        Código FIPE
                      </Typography>
                      <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        {fipeDetail.codeFipe}
                      </Typography>
                    </Box>
                  </Stack>
                </Grid>

                <Grid size={{ xs: 12, sm: 4 }}>
                  <Stack
                    direction="row"
                    spacing={1}
                    sx={{ alignItems: "center" }}
                  >
                    <MonetizationOnRoundedIcon
                      fontSize="small"
                      sx={{ color: "success.main" }}
                    />
                    <Box>
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ display: "block" }}
                      >
                        Preço Médio de Mercado
                      </Typography>
                      <Typography
                        variant="body2"
                        sx={{ fontWeight: 700, color: "success.main" }}
                      >
                        {fipeDetail.price}
                      </Typography>
                    </Box>
                  </Stack>
                </Grid>

                <Grid size={{ xs: 12, sm: 4 }}>
                  <Stack
                    direction="row"
                    spacing={1}
                    sx={{ alignItems: "center" }}
                  >
                    <LocalGasStationRoundedIcon
                      fontSize="small"
                      sx={{ color: "text.secondary" }}
                    />
                    <Box>
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ display: "block" }}
                      >
                        Combustível
                      </Typography>
                      <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        {fipeDetail.fuel}
                      </Typography>
                    </Box>
                  </Stack>
                </Grid>
              </Grid>
            </Paper>
          )}

          {/* Dados Finais do Veículo */}
          <Box>
            <Typography
              variant="caption"
              sx={{
                fontWeight: 600,
                color: "text.secondary",
                textTransform: "uppercase",
                letterSpacing: 0.5,
                display: "block",
                mb: 1.5,
              }}
            >
              3. Identificação e Quilometragem
            </Typography>

            <Grid container spacing={2}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Apelido / Identificador"
                  placeholder="Ex: Meu Carro, Carro da Família"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                  disabled={loading}
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Km Atual"
                  type="number"
                  placeholder="0"
                  value={lastMileage}
                  onChange={(e) => setLastMileage(Number(e.target.value))}
                  disabled={loading}
                  slotProps={{
                    htmlInput: {
                      min: 0,
                    },
                  }}
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Ano de Fabricação"
                  type="number"
                  placeholder="Ex: 2022"
                  value={yearManufacture || ""}
                  onChange={(e) => setYearManufacture(Number(e.target.value))}
                  disabled={loading}
                  slotProps={{
                    htmlInput: {
                      min: 1900,
                      max: currentYear + 1,
                    },
                  }}
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Ano do Modelo"
                  type="number"
                  placeholder="Ex: 2023"
                  value={yearModel || ""}
                  onChange={(e) => setYearModel(Number(e.target.value))}
                  disabled={loading}
                  slotProps={{
                    htmlInput: {
                      min: 1900,
                      max: currentYear + 2,
                    },
                  }}
                />
              </Grid>
            </Grid>
          </Box>

          {/* Upload de Foto do Veículo */}
          <ImageUploadField
            value={imageUrl}
            onChange={setImageUrl}
            vehicleType={vehicleType}
            disabled={loading}
          />
        </Stack>
      </DialogContent>

      <DialogActions sx={{ p: 2.5, gap: 1 }}>
        <Button
          onClick={handleClose}
          variant="outlined"
          color="inherit"
          disabled={loading}
        >
          Cancelar
        </Button>
        <Button
          onClick={() => handleSubmit()}
          variant="contained"
          disabled={
            loading ||
            user?.emailVerified === false ||
            !selectedYear ||
            !name.trim()
          }
        >
          {loading ? (
            <CircularProgress size={22} sx={{ color: "#FFFFFF" }} />
          ) : (
            "Cadastrar Veículo"
          )}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
