#include <iostream>
#include <fstream>
#include <set>
#include <iomanip>
#include "pin.H"

// Output file stream
static std::ofstream CoverageFile;

// Base address of the main executable (set when the main image loads)
static ADDRINT imageBase = 0;

// A set of RVAs that have been executed
static std::set<ADDRINT> coveredRvas;

//------------------------------------------------------------------------------
// Function: RecordBasicBlock
//   Called at runtime whenever a basic block executes. 
//   We calculate an RVA by subtracting the base address and insert it into 'coveredRvas'.
//------------------------------------------------------------------------------
static VOID RecordBasicBlock(ADDRINT addr)
{
    if (imageBase != 0)
    {
        ADDRINT rva = addr - imageBase;
        coveredRvas.insert(rva);
    }
}

//------------------------------------------------------------------------------
// Function: OnImageLoad
//   Called once for each loaded image. We detect the main executable, 
//   record its load address, and store it in 'imageBase'.
//------------------------------------------------------------------------------
static VOID OnImageLoad(IMG img, VOID* v)
{
    if (IMG_IsMainExecutable(img))
    {
        imageBase = IMG_LowAddress(img);
        std::cerr << "[+] Main Executable loaded at 0x" 
                  << std::hex << imageBase << std::endl;
    }
}

//------------------------------------------------------------------------------
// Function: OnTrace
//   Called for every TRACE in the program. We iterate over the basic blocks (BBL)
//   in each TRACE and instrument them by inserting a call to 'RecordBasicBlock'.
//------------------------------------------------------------------------------
static VOID OnTrace(TRACE trace, VOID* v)
{
    for (BBL bbl = TRACE_BblHead(trace); BBL_Valid(bbl); bbl = BBL_Next(bbl))
    {
        BBL_InsertCall(bbl,
                       IPOINT_ANYWHERE,
                       (AFUNPTR)RecordBasicBlock,
                       IARG_ADDRINT, BBL_Address(bbl),
                       IARG_END);
    }
}

//------------------------------------------------------------------------------
// Function: OnFini
//   Called when the application finishes. We dump all recorded RVAs to a file.
//------------------------------------------------------------------------------
static VOID OnFini(INT32 code, VOID* v)
{
    CoverageFile.open("coverage_rva.out");
    if (!CoverageFile.is_open())
    {
        std::cerr << "[!] Could not open coverage_rva.out for writing!\n";
        return;
    }

    for (auto rva : coveredRvas)
    {
        CoverageFile << "CCu [COVERED] @0x" << std::hex << rva << std::endl;
    }

    CoverageFile.close();
    std::cerr << "[+] Code coverage saved to coverage_rva.out" << std::endl;
}

//------------------------------------------------------------------------------
// Function: main
//   PIN tool entry point. Initializes symbol processing, registers callbacks, 
//   and starts the target program under PIN instrumentation.
//------------------------------------------------------------------------------
int main(int argc, char* argv[])
{
    // Initialize symbol processing (needed for IMG routines, etc.)
    PIN_InitSymbols();

    // Initialize PIN
    if (PIN_Init(argc, argv))
    {
        std::cerr << "[!] PIN Initialization failed." << std::endl;
        return 1;
    }

    // Register instrumentation callbacks
    IMG_AddInstrumentFunction(OnImageLoad, nullptr);
    TRACE_AddInstrumentFunction(OnTrace, nullptr);
    PIN_AddFiniFunction(OnFini, nullptr);

    // Start the program (never returns)
    PIN_StartProgram();
    return 0; 
}

