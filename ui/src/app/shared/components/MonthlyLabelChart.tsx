import * as React from 'react';
import { Card, CardBody, Title, Spinner } from '@patternfly/react-core';
import {
  Chart,
  ChartAxis,
  ChartLine,
  ChartVoronoiContainer,
  ChartGroup,
} from '@patternfly/react-charts/victory';
import { IMonthlyTagReport } from '@app/shared/models/label';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface MonthlyLabelChartProps {
  data: IMonthlyTagReport[];
  loading?: boolean;
  title?: string;
  labelNames?: string[];
  startDate?: string;
  endDate?: string;
}

const MonthlyLabelChart: React.FC<MonthlyLabelChartProps> = ({
  data,
  loading = false,
  title = 'Monthly Label Amounts',
  labelNames = [],
  startDate,
  endDate,
}) => {
  const { theme } = useTheme();
  // Validate date range spans multiple months
  const dateRangeError = React.useMemo(() => {
    if (!startDate || !endDate) return null;

    const start = new Date(startDate);
    const end = new Date(endDate);

    // Check if dates span multiple months
    const startMonthYear = `${start.getFullYear()}-${start.getMonth()}`;
    const endMonthYear = `${end.getFullYear()}-${end.getMonth()}`;

    if (startMonthYear === endMonthYear) {
      return 'Date range must span multiple months to display monthly trends';
    }

    return null;
  }, [startDate, endDate]);

  const chartData = React.useMemo(() => {
    if (!data || data.length === 0 || labelNames.length === 0) return {};
    if (dateRangeError) return {};

    // Filter data by date range if provided
    let filteredData = data;
    if (startDate && endDate) {
      const start = new Date(startDate);
      const end = new Date(endDate);

      filteredData = data.filter((item) => {
        const itemDate = new Date(item.year, item.month - 1); // month is 1-based in data
        return itemDate >= start && itemDate <= end;
      });
    }

    // Group data by label
    const labelGroups: { [labelName: string]: Array<{ name: string; x: string; y: number }> } = {};

    labelNames.forEach((labelName) => {
      const labelData = filteredData
        .filter((item) => item.tag === labelName) // Keep as 'tag' for backward compatibility
        .sort((a, b) => {
          if (a.year !== b.year) return a.year - b.year;
          return a.month - b.month;
        })
        .map((item) => ({
          x: `${item.year}-${String(item.month).padStart(2, '0')}`,
          y: item.amount,
          name: labelName,
        }));

      if (labelData.length > 0) {
        labelGroups[labelName] = labelData;
      }
    });

    return labelGroups;
  }, [data, labelNames, startDate, endDate, dateRangeError]);

  // Calculate the maximum Y value from all chart data for dynamic domain
  const maxYValue = React.useMemo(() => {
    const allValues: number[] = [];
    Object.values(chartData).forEach((labelData) => {
      labelData.forEach((point) => {
        allValues.push(point.y);
      });
    });

    if (allValues.length === 0) return 1000; // Default fallback

    const max = Math.max(...allValues);
    // Add 10% padding to the max value for better visualization
    return Math.ceil(max * 1.1);
  }, [chartData]);

  // Extract all unique x values from chartData for x-axis
  const xAxisValues = React.useMemo(() => {
    const allXValues = new Set<string>();
    Object.values(chartData).forEach((labelData) => {
      labelData.forEach((point) => {
        allXValues.add(point.x);
      });
    });

    // Sort the x values chronologically
    return Array.from(allXValues).sort();
  }, [chartData]);

  // Create legend data from selected labels
  const legendData = React.useMemo(() => {
    return labelNames.map((labelName) => {
      // Extract key part from "key=value" format, or use full string if no '=' found
      const labelKey = labelName.includes(':') ? labelName.split('=')[0] : labelName;
      return {
        name: labelKey,
      };
    });
  }, [labelNames]);

  if (loading) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <Spinner size="lg" />
          </div>
        </CardBody>
      </Card>
    );
  }

  // Show error message if date range is invalid
  if (dateRangeError) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <div style={{ textAlign: 'center', color: 'var(--pf-v6-global--danger-color--100)' }}>
              <strong>Error:</strong> {dateRangeError}
            </div>
          </div>
        </CardBody>
      </Card>
    );
  }

  if (Object.keys(chartData).length === 0) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <div>No data available</div>
          </div>
        </CardBody>
      </Card>
    );
  }

  return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
                                       <div style={{
            height: '600px',
            width: '100%'
          }}>
        <Chart
          ariaDesc="Monthly label amounts chart"
          ariaTitle={title}
          containerComponent={
            <ChartVoronoiContainer
              labels={({ datum }) => `${datum.name}: ${datum.y.toLocaleString('fr-FR', {
                style: 'currency',
                currency: 'EUR',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
              })}`}
              constrainToVisibleArea
            />
          }
          legendData={legendData}
          legendOrientation="vertical"
          legendPosition="right"
          maxDomain={{ y: maxYValue }}
          minDomain={{ y: 0 }}
          name="monthlyLabelChart"
          padding={{
            bottom: 80,
            left: 80,
            right: 120,
            top: 50,
          }}
        >
         <ChartAxis
           tickValues={xAxisValues}
           tickFormat={(value) => {
             // Format month-year string for display (e.g., "2024-01" -> "Jan 2024")
             const [year, month] = value.split('-');
             const monthNames = [
               'Jan',
               'Feb',
               'Mar',
               'Apr',
               'May',
               'Jun',
               'Jul',
               'Aug',
               'Sep',
               'Oct',
               'Nov',
               'Dec',
             ];
             return `${monthNames[parseInt(month) - 1]} ${year}`;
           }}
         />
                   <ChartAxis dependentAxis showGrid tickFormat={(value) =>
            value.toLocaleString('fr-FR', {
              style: 'currency',
              currency: 'EUR',
              minimumFractionDigits: 0,
              maximumFractionDigits: 0,
            })
          } />
         <ChartGroup>
           {Object.entries(chartData).map(([labelName, labelData]) => (
             <ChartLine key={labelName} data={labelData} />
           ))}
                  </ChartGroup>
       </Chart>
     </div>
     {labelNames.length > 0 && (
       <div
         style={{
           textAlign: 'center',
           marginTop: '1rem',
           fontSize: '0.875rem',
           color: 'var(--pf-v6-global--Color--200)',
         }}
       >
         Showing amounts for: <strong>{labelNames.join(', ')}</strong>
       </div>
     )}
     </CardBody>
     </Card>
  );
};

export { MonthlyLabelChart };
