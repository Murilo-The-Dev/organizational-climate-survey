
import Sidebar from "../layout/Sidebar";

type DashboardLayoutProps = {
    children: React.ReactNode;
};

const DashboardLayout = ({ children }: DashboardLayoutProps) => {
    return (
        <div className="flex w-full min-h-screen bg-muted/40">
            <Sidebar isOpen={true} />
            <main className="flex-1 p-4 md:p-8">
                {children}
            </main>
        </div>
    );
};

export default DashboardLayout;