import React from "react";
import { Box, BoxProps } from "@mui/material";

export function MotorcycleVector(props: BoxProps) {
  return (
    <Box
      component="svg"
      viewBox="0 0 800 450"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      sx={{
        width: "100%",
        height: "100%",
        display: "block",
        ...props.sx,
      }}
      {...props}
    >
      <defs>
        {/* Background Glow */}
        <radialGradient id="motoBgGlow" cx="50%" cy="45%" r="60%">
          <stop offset="0%" stopColor="#38BDF8" stopOpacity="0.15" />
          <stop offset="60%" stopColor="#0284C7" stopOpacity="0.05" />
          <stop offset="100%" stopColor="#0F172A" stopOpacity="0" />
        </radialGradient>

        {/* Motorcycle Tank & Fairing Metallic Gradient */}
        <linearGradient
          id="motoBodyGrad"
          x1="280"
          y1="140"
          x2="620"
          y2="300"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#38BDF8" />
          <stop offset="35%" stopColor="#0284C7" />
          <stop offset="80%" stopColor="#0369A1" />
          <stop offset="100%" stopColor="#0F172A" />
        </linearGradient>

        {/* Metallic Frame & Engine */}
        <linearGradient
          id="engineGrad"
          x1="300"
          y1="240"
          x2="520"
          y2="360"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#475569" />
          <stop offset="50%" stopColor="#334155" />
          <stop offset="100%" stopColor="#0F172A" />
        </linearGradient>

        {/* Exhaust Chrome */}
        <linearGradient
          id="chromeGrad"
          x1="200"
          y1="320"
          x2="480"
          y2="360"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#F8FAFC" />
          <stop offset="50%" stopColor="#94A3B8" />
          <stop offset="100%" stopColor="#475569" />
        </linearGradient>

        {/* Rim Radial */}
        <radialGradient id="motoRimGrad" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#E2E8F0" />
          <stop offset="70%" stopColor="#64748B" />
          <stop offset="100%" stopColor="#0F172A" />
        </radialGradient>

        {/* Headlight Glow */}
        <linearGradient
          id="motoHeadlightGlow"
          x1="620"
          y1="180"
          x2="760"
          y2="210"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#38BDF8" stopOpacity="0.85" />
          <stop offset="100%" stopColor="#38BDF8" stopOpacity="0" />
        </linearGradient>

        {/* Ground Shadow */}
        <radialGradient id="motoShadow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#0F172A" stopOpacity="0.6" />
          <stop offset="70%" stopColor="#0F172A" stopOpacity="0.2" />
          <stop offset="100%" stopColor="#0F172A" stopOpacity="0" />
        </radialGradient>
      </defs>

      {/* Background Glow */}
      <rect width="800" height="450" fill="url(#motoBgGlow)" />

      {/* Ground Contact Shadows */}
      <ellipse cx="430" cy="375" rx="300" ry="24" fill="url(#motoShadow)" />
      <ellipse cx="260" cy="370" rx="70" ry="16" fill="#000000" opacity="0.6" />
      <ellipse cx="585" cy="370" rx="75" ry="16" fill="#000000" opacity="0.6" />

      {/* Headlight Beam */}
      <polygon
        points="620,180 770,170 780,260 630,230"
        fill="url(#motoHeadlightGlow)"
      />

      {/* EXHAUST PIPE (Chrome) */}
      <path
        d="M380 290 
           C400 320, 390 350, 320 350 
           L220 340 
           C200 338, 190 325, 205 315 
           L290 320 
           C340 325, 360 305, 370 290 Z"
        fill="url(#chromeGrad)"
        stroke="#475569"
        strokeWidth="2"
      />
      <ellipse
        cx="205"
        cy="325"
        rx="10"
        ry="14"
        fill="#0F172A"
        stroke="#CBD5E1"
        strokeWidth="2"
      />

      {/* MAIN FRAME & ENGINE BLOCK */}
      <path
        d="M340 230 
           L470 200 
           L530 250 
           L490 340 
           L360 340 
           L320 270 Z"
        fill="url(#engineGrad)"
        stroke="#1E293B"
        strokeWidth="3"
      />

      {/* Engine Cooling Fins & Crankcase Details */}
      <rect x="375" y="260" width="80" height="8" rx="4" fill="#64748B" />
      <rect x="370" y="275" width="90" height="8" rx="4" fill="#64748B" />
      <rect x="375" y="290" width="85" height="8" rx="4" fill="#64748B" />
      <rect x="380" y="305" width="75" height="8" rx="4" fill="#64748B" />
      <circle
        cx="430"
        cy="320"
        r="22"
        fill="#1E293B"
        stroke="#64748B"
        strokeWidth="3"
      />
      <circle cx="430" cy="320" r="8" fill="#38BDF8" />

      {/* REAR SWINGARM & DRIVE CHAIN */}
      <polygon
        points="260,345 380,310 375,335 260,360"
        fill="#334155"
        stroke="#1E293B"
        strokeWidth="2"
      />
      <circle
        cx="260"
        cy="345"
        r="24"
        fill="#1E293B"
        stroke="#64748B"
        strokeWidth="2"
      />

      {/* MONOSHOCK SUSPENSION (Yellow / Sport Spring) */}
      <line
        x1="330"
        y1="280"
        x2="380"
        y2="320"
        stroke="#F59E0B"
        strokeWidth="8"
        strokeLinecap="round"
      />
      <line
        x1="330"
        y1="280"
        x2="380"
        y2="320"
        stroke="#0F172A"
        strokeWidth="2"
        strokeDasharray="4 4"
      />

      {/* SEAT & TAIL SECTION (3/4 Sport Styling) */}
      <path
        d="M200 240 
           C250 245, 300 250, 360 250 
           L440 240 
           C410 230, 370 215, 320 220 
           L240 225 
           C210 228, 190 232, 200 240 Z"
        fill="#0F172A"
        stroke="#334155"
        strokeWidth="2"
      />
      {/* Passenger Seat / Tail Cowl */}
      <path
        d="M190 230 
           L280 225 
           L290 210 
           L210 215 Z"
        fill="url(#motoBodyGrad)"
      />
      {/* Taillight */}
      <polygon
        points="185,228 195,224 198,235 188,237"
        fill="#EF4444"
        stroke="#F87171"
        strokeWidth="1"
      />

      {/* FUEL TANK (Sculpted Modern Street Tank) */}
      <path
        d="M340 240 
           C360 200, 420 160, 490 160 
           C535 160, 560 185, 565 215 
           L530 245 
           C480 255, 410 255, 340 240 Z"
        fill="url(#motoBodyGrad)"
        stroke="#38BDF8"
        strokeWidth="2"
      />
      {/* Tank Highlight Streak */}
      <path
        d="M375 220 C420 180, 470 170, 520 175 C490 195, 430 210, 375 220 Z"
        fill="#FFFFFF"
        opacity="0.35"
      />

      {/* FRONT FORK & HANDLEBARS */}
      {/* Handlebar & Grips */}
      <path
        d="M525 145 L565 140 L595 155"
        stroke="#CBD5E1"
        strokeWidth="6"
        strokeLinecap="round"
      />
      <rect
        x="585"
        y="148"
        width="18"
        height="8"
        rx="3"
        fill="#0F172A"
        transform="rotate(25 585 148)"
      />
      {/* Rearview Mirror */}
      <line
        x1="560"
        y1="140"
        x2="550"
        y2="105"
        stroke="#94A3B8"
        strokeWidth="3"
        strokeLinecap="round"
      />
      <ellipse
        cx="545"
        cy="100"
        rx="14"
        ry="9"
        fill="#0F172A"
        stroke="#38BDF8"
        strokeWidth="1.5"
      />

      {/* Triple Tree & Front Fork Tubes */}
      <line
        x1="560"
        y1="165"
        x2="615"
        y2="345"
        stroke="#38BDF8"
        strokeWidth="8"
        strokeLinecap="round"
      />
      <line
        x1="550"
        y1="170"
        x2="600"
        y2="340"
        stroke="#94A3B8"
        strokeWidth="6"
        strokeLinecap="round"
      />

      {/* Front Headlight Cowl / Mask */}
      <path
        d="M570 165 
           C600 170, 625 185, 630 205 
           C625 225, 595 235, 570 230 Z"
        fill="url(#motoBodyGrad)"
        stroke="#0284C7"
        strokeWidth="2"
      />
      {/* LED Projector Lens */}
      <polygon
        points="605,185 628,198 622,215 600,205"
        fill="#E0F2FE"
        stroke="#38BDF8"
        strokeWidth="2"
      />
      <circle cx="615" cy="200" r="5" fill="#38BDF8" />

      {/* Front Fender / Mudguard */}
      <path
        d="M560 270 C590 260, 640 280, 660 310 L630 315 C615 295, 580 285, 560 285 Z"
        fill="url(#motoBodyGrad)"
      />

      {/* REAR WHEEL */}
      <g>
        <circle
          cx="260"
          cy="345"
          r="62"
          fill="#1E293B"
          stroke="#0F172A"
          strokeWidth="5"
        />
        <circle cx="260" cy="345" r="46" fill="url(#motoRimGrad)" />
        {/* Y-Spoke Wheels */}
        <circle cx="260" cy="345" r="16" fill="#0F172A" />
        <line
          x1="260"
          y1="305"
          x2="260"
          y2="385"
          stroke="#F8FAFC"
          strokeWidth="4"
          strokeLinecap="round"
        />
        <line
          x1="225"
          y1="345"
          x2="295"
          y2="345"
          stroke="#F8FAFC"
          strokeWidth="4"
          strokeLinecap="round"
        />
        <line
          x1="235"
          y1="320"
          x2="285"
          y2="370"
          stroke="#CBD5E1"
          strokeWidth="3"
          strokeLinecap="round"
        />
        <line
          x1="235"
          y1="370"
          x2="285"
          y2="320"
          stroke="#CBD5E1"
          strokeWidth="3"
          strokeLinecap="round"
        />
        {/* Brake Disc */}
        <circle
          cx="260"
          cy="345"
          r="28"
          fill="none"
          stroke="#94A3B8"
          strokeWidth="4"
          strokeDasharray="4 3"
        />
        <circle cx="260" cy="345" r="6" fill="#38BDF8" />
      </g>

      {/* FRONT WHEEL (Larger in 3/4 perspective) */}
      <g>
        <circle
          cx="610"
          cy="345"
          r="68"
          fill="#1E293B"
          stroke="#0F172A"
          strokeWidth="5"
        />
        <circle cx="610" cy="345" r="50" fill="url(#motoRimGrad)" />
        {/* Y-Spoke Wheels */}
        <circle cx="610" cy="345" r="18" fill="#0F172A" />
        <line
          x1="610"
          y1="300"
          x2="610"
          y2="390"
          stroke="#F8FAFC"
          strokeWidth="4.5"
          strokeLinecap="round"
        />
        <line
          x1="570"
          y1="345"
          x2="650"
          y2="345"
          stroke="#F8FAFC"
          strokeWidth="4.5"
          strokeLinecap="round"
        />
        <line
          x1="580"
          y1="315"
          x2="640"
          y2="375"
          stroke="#CBD5E1"
          strokeWidth="3.5"
          strokeLinecap="round"
        />
        <line
          x1="580"
          y1="375"
          x2="640"
          y2="315"
          stroke="#CBD5E1"
          strokeWidth="3.5"
          strokeLinecap="round"
        />
        {/* Perforated Disc Brake */}
        <circle
          cx="610"
          cy="345"
          r="34"
          fill="none"
          stroke="#CBD5E1"
          strokeWidth="5"
          strokeDasharray="5 3"
        />
        {/* Sport Caliper (Gold/Cyan) */}
        <path
          d="M580 330 C580 320, 595 315, 605 320 L600 340 Z"
          fill="#0284C7"
          stroke="#38BDF8"
          strokeWidth="1.5"
        />
        <circle cx="610" cy="345" r="7" fill="#38BDF8" />
      </g>
    </Box>
  );
}
